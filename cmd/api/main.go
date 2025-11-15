package main

import (
	"context"
	"devconnectrelations/internal/adapters/inbound/consumer/profile_created"
	city_rest "devconnectrelations/internal/adapters/inbound/rest/city"
	cityrelation "devconnectrelations/internal/adapters/inbound/rest/city_relation"
	rest "devconnectrelations/internal/adapters/inbound/rest/profile"
	relation_controller "devconnectrelations/internal/adapters/inbound/rest/relation"
	stack_rest "devconnectrelations/internal/adapters/inbound/rest/stack"
	stack_relation_rest "devconnectrelations/internal/adapters/inbound/rest/stack_relation"
	"devconnectrelations/internal/adapters/outbound/repository"
	cityRepository "devconnectrelations/internal/adapters/outbound/repository/city"
	profileRepository "devconnectrelations/internal/adapters/outbound/repository/profile"
	relationRepository "devconnectrelations/internal/adapters/outbound/repository/profile_relation"
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
	repo := profileRepository.NewNeo4jProfileRepository(driver)
	profile_service := service.CreateNewProfileService(repo)
	profile_controller := rest.CreateNewProfileController(*profile_service)
	router.POST("/profile", profile_controller.CreateProfile)
	router.DELETE("/profile/:id", profile_controller.DeleteProfile)
	router.GET("/profile/:id", profile_controller.GetProfileByID)
	return profile_service
}

func setRelation(router *gin.Engine, driver neo4j.DriverWithContext, profile_service *service.ProfileService) {
	repo := relationRepository.NewNeo4jRelationRepository(driver)
	relation_service := service.CreateRelationService(repo)
	relation_controller := relation_controller.CreateNewRelationsController(*relation_service)
	router.POST("/relation", relation_controller.CreateRelation)
	router.GET("/relation/:fromId", relation_controller.GetAllRelationsByFromId)
	router.PATCH("/relation/accept/:fromId/:toId", relation_controller.AcceptRelation)
	router.GET("/relation/pending/:fromId", relation_controller.GetAllRelationPendingByFromId)
}

func setStack(router *gin.Engine, driver neo4j.DriverWithContext) *service.StackService {
	repo := repository.NewNeo4jStackRepository(driver)
	stack_service := service.CreateStackService(repo)
	stack_controller := stack_rest.CreateNewStackController(*stack_service)
	router.POST("/stack", stack_controller.CreateStack)
	router.GET("/stack/:name", stack_controller.GetStackByName)
	router.DELETE("/stack/:name", stack_controller.DeleteStack)
	return stack_service
}

func setStackRelation(router *gin.Engine, driver neo4j.DriverWithContext) *service.StackRelationService {
	repo := relationRepository.NewNeo4jStackRelationRepository(driver)
	stack_relation_service := service.CreateStackRelationService(repo)
	stack_relation_controller := stack_relation_rest.CreateNewStackRelationController(stack_relation_service)
	router.POST("/stack-relation", stack_relation_controller.CreateStackRelation)
	router.DELETE("/stack-relation", stack_relation_controller.DeleteStackRelation)
	return stack_relation_service
}

func setCity(router *gin.Engine, driver neo4j.DriverWithContext) *service.CityService {
	repo := cityRepository.NewNeo4jCityRepository(driver)
	city_service := service.NewCityService(repo)
	city_controller := city_rest.CreateNewCityController(*city_service)
	router.POST("/city", city_controller.CreateCity)
	router.GET("/city/:fullName", city_controller.GetCityByFullName)
	return city_service
}

func setCityRelation(router *gin.Engine, driver neo4j.DriverWithContext, cityService *service.CityService) *service.CityRelationService {
	repo := relationRepository.NewNeo4jRelationCityRepository(&driver)
	city_relation_service := service.CreateNewCityRelationService(repo, cityService)
	city_relation_controller := cityrelation.CreateNewCityRelationController(*city_relation_service)
	router.POST("/city-relation", city_relation_controller.CreateCityRelation)
	return city_relation_service
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
	stackService := setStack(router, driver)
	stackRelationService := setStackRelation(router, driver)
	cityService := setCity(router, driver)
	cityRelationService := setCityRelation(router, driver, cityService)
	kafka_brokers := []string{GetEnv("KAFKA_SERVER", "localhost:9092")}
	kafka_profile_create_topic := GetEnv("KAFKA_PROFILE_CREATED_TOPIC", "dev-profile.created.v1")
	kafka_group_id := GetEnv("KAFKA_GROUP_ID", "dev-connect-relations-group")
	consumer := profile_created.NewKafkaProfileCreatedConsumer(kafka_brokers, kafka_profile_create_topic, kafka_group_id, profile_service, stackService, stackRelationService, cityService, cityRelationService)
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
