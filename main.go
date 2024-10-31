package main

import (
	"websocket-chat/agent"
	"websocket-chat/handler"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Initialize Logrus settings
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	log.SetLevel(log.DebugLevel)

	// Initialize Agent
	masterAgent := agent.NewMasterAgent()
	supportAgent := agent.NewSupportAgent()

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
