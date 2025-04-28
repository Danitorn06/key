package handlers

import (
    "myproject/database"
    "myproject/models"
    "net/http"
    "database/sql"

    "github.com/gin-gonic/gin"
)

func RegisterWebsite(c *gin.Context) {
    var website models.Website
    if err := c.ShouldBindJSON(&website); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    claimsInterface, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    claims := claimsInterface.(map[string]interface{})
    username := claims["preferred_username"].(string)

    var affiliatorID int
    err := database.DB.QueryRow("SELECT id FROM clients WHERE name = $1", username).Scan(&affiliatorID)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Affiliator not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }

    _, err = database.DB.Exec(
        "INSERT INTO affiliator_websites (affiliator_id, website_url) VALUES ($1, $2)",
        affiliatorID, website.WebsiteURL,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register website"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Website registered successfully"})
}
