package consumer

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/service"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaProfileCreatedConsumer struct {
	reader               *kafka.Reader
	service              *service.ProfileService
	stackService         *service.StackService
	stackRelationService *service.StackRelationService
}

func NewKafkaProfileCreatedConsumer(brokers []string, topic, groupID string, service *service.ProfileService, stackService *service.StackService, stackRelationService *service.StackRelationService) *KafkaProfileCreatedConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaProfileCreatedConsumer{
		reader:               reader,
		service:              service,
		stackService:         stackService,
		stackRelationService: stackRelationService,
	}
}

func (c *KafkaProfileCreatedConsumer) Consume(ctx context.Context) error {
	fmt.Println("🚀 Kafka consumer started for topic:", c.reader.Config().Topic)

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var profile entities.Profile
		var createdEvent CreatedProfileEvent
		if err := json.Unmarshal(m.Value, &createdEvent); err != nil {
			fmt.Println("❌ Erro ao deserializar mensagem:", err)
			continue
		}
		profile.ConnectId = createdEvent.Id
		profile.Name = createdEvent.Name
		fmt.Printf("📩 Mensagem recebida: %+v\n", profile)

		if _, err := c.service.CreateProfile(ctx, &profile); err != nil {
			fmt.Println("❌ Erro ao criar perfil:", err)
			continue
		}
		var stacks []*entities.Stack
		for _, stackName := range createdEvent.Stack {
			existsStack, getErr := c.stackService.GetStackByName(ctx, stackName)
			if getErr != nil {
				fmt.Println("❌ Erro ao verificar stack existente:", getErr)
				continue
			}
			if existsStack == (entities.Stack{}) {
				newStack, createErr := c.stackService.CreateStack(ctx, stackName)
				if createErr != nil {
					fmt.Println("❌ Erro ao criar stack:", createErr)
					continue
				}
				stacks = append(stacks, &newStack)
			} else {
				stacks = append(stacks, &existsStack)
			}
			if _, relErr := c.stackRelationService.CreateStackRelation(ctx, stackName, profile.ConnectId); relErr != nil {
				fmt.Println("❌ Erro ao criar relação stack-profile:", relErr)
				continue
			}
		}
		fmt.Println("✅ Perfil criado no Neo4j com sucesso! ID:", profile.ConnectId)
		for _, stack := range stacks {
			fmt.Println("✅ Stack criada no Neo4j com sucesso! Nome:", stack.Name)
			fmt.Printf("✅ Relação entre Profile ID %d e Stack %s criada com sucesso!\n", profile.ConnectId, stack.Name)
		}

		return nil
	}

}
