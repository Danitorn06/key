package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"myproject/database"
	"myproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func RegisterClient(c *gin.Context) {
	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	client.APIKey = generateAPIKey()

	_, err := database.DB.Exec("INSERT INTO clients (name, email, api_key) VALUES ($1, $2, $3)", client.Name, client.Email, client.APIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client registered", "api_key": client.APIKey})
}
