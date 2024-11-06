package main

import (
	"websocket-chat/agent"
	"websocket-chat/database"
	"websocket-chat/handler"
	"websocket-chat/repository"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Initialize Logrus settings
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	logger.SetLevel(log.DebugLevel)

	// Initialize DB
	database := database.NewMongoDB(logger)
	database.ConnectToMongoDB()

	// Initialize Repository
	repository := repository.NewRepository(database.GetMongoClient(), logger)

	// Initialize Agent
	masterAgent := agent.NewMasterAgent(logger, repository)
	supportAgent := agent.NewSupportAgent(logger, repository)

	// Initialize Handler
	h := handler.NewHandler(masterAgent, supportAgent)

	// Initialize Gin router
	router := gin.Default()
	router.Use(gin.Recovery())

	router.GET("/ws/client", h.ClientWs)
	router.GET("/ws/support", h.SupportAgentWs)

	// Start server
	port := ":8080"
	log.Infof("Server started on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
