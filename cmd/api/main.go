package main

import (
	"context"
	rest "devconnectrelations/internal/adapters/inbound/rest/profile"
	relation_controller "devconnectrelations/internal/adapters/inbound/rest/relation"
	"devconnectrelations/internal/adapters/outbound/repository"
	"devconnectrelations/internal/domain/service"
	"fmt"
	"os"

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

func main() {
	ctx := context.Background()
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
	repo := repository.NewNeo4jProfileRepository(driver)
	profile_service := service.CreateNewProfileService(repo)
	profile_controller := rest.CreateNewProfileController(*profile_service)
	realtion_repo := repository.NewNeo4jRelationRepository(driver)
	relation_service := service.CreateRelationService(realtion_repo)
	rest_relation_controller := relation_controller.CreateNewRelationsController(*relation_service)
	router.POST("/relation", rest_relation_controller.CreateRelation)
	router.POST("/profile", profile_controller.CreateProfile)
	router.GET("/relation/:fromId", rest_relation_controller.GetAllRelationsByFromId)
	router.DELETE("/profile/:id", profile_controller.DeleteProfile)
	router.Run()
}
