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
	reader  *kafka.Reader
	service *service.ProfileService
}

func NewKafkaProfileCreatedConsumer(brokers []string, topic, groupID string, service *service.ProfileService) *KafkaProfileCreatedConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaProfileCreatedConsumer{
		reader:  reader,
		service: service,
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

		fmt.Println("✅ Perfil criado no Neo4j com sucesso! ID:", profile.ConnectId)
		return nil
	}

}
