package main

import (
	"context"
	"devconnectrelations/internal/adapters/inbound/consumer"
	rest "devconnectrelations/internal/adapters/inbound/rest/profile"
	relation_controller "devconnectrelations/internal/adapters/inbound/rest/relation"
	stack_rest "devconnectrelations/internal/adapters/inbound/rest/stack"
	"devconnectrelations/internal/adapters/outbound/repository"
	"devconnectrelations/internal/domain/service"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Println(value)
		return value
	}
	fmt.Println(fallback)
	return fallback
}

func setProfile(router *gin.Engine, driver neo4j.DriverWithContext) *service.ProfileService {
	repo := repository.NewNeo4jProfileRepository(driver)
	profile_service := service.CreateNewProfileService(repo)
	profile_controller := rest.CreateNewProfileController(*profile_service)
	router.POST("/profile", profile_controller.CreateProfile)
	router.DELETE("/profile/:id", profile_controller.DeleteProfile)
	return profile_service
}

func setRelation(router *gin.Engine, driver neo4j.DriverWithContext, profile_service *service.ProfileService) {
	repo := repository.NewNeo4jRelationRepository(driver)
	relation_service := service.CreateRelationService(repo)
	relation_controller := relation_controller.CreateNewRelationsController(*relation_service)
	router.POST("/relation", relation_controller.CreateRelation)
	router.GET("/relation/:fromId", relation_controller.GetAllRelationsByFromId)
	router.PATCH("/relation/accept/:fromId/:toId", relation_controller.AcceptRelation)
	router.GET("/relation/pending/:fromId", relation_controller.GetAllRelationPendingByFromId)
}

func setStack(router *gin.Engine, driver neo4j.DriverWithContext) {
	repo := repository.NewNeo4jStackRepository(driver)
	stack_service := service.CreateStackService(repo)
	stack_controller := stack_rest.CreateNewStackController(*stack_service)
	router.POST("/stack", stack_controller.CreateStack)
	router.GET("/stack/:name", stack_controller.GetStackByName)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	router := gin.Default()
	dbUri := GetEnv("NEO4J_URI", "neo4j://localhost:7687")
	dbUser := GetEnv("NEO4J_USER", "neo4j")
	dbPassword := GetEnv("NEO4J_PASSWORD", "Kadeira4.0")
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	profile_service := setProfile(router, driver)
	setRelation(router, driver, profile_service)
	setStack(router, driver)
	kafka_brokers := []string{GetEnv("KAFKA_SERVER", "localhost:9092")}
	kafka_profile_create_topic := GetEnv("KAFKA_PROFILE_CREATED_TOPIC", "dev-profile.created.v1")
	kafka_group_id := GetEnv("KAFKA_GROUP_ID", "dev-connect-relations-group")
	consumer := consumer.NewKafkaProfileCreatedConsumer(kafka_brokers, kafka_profile_create_topic, kafka_group_id, profile_service)
	go func() {
		if err := consumer.Consume(ctx); err != nil {
			fmt.Println("kafka consumer error:", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("🔔 shutdown signal received")
		cancel()
		os.Exit(0)
	}()
	router.Run(":8082")
}
