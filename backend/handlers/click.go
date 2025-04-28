package handlers

import (
    "myproject/database"
    "myproject/models"
    "net/http"

    "github.com/gin-gonic/gin"
)

func LogClick(c *gin.Context) {
    var clickLog models.ClickLog
    if err := c.ShouldBindJSON(&clickLog); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    claimsInterface, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    claims := claimsInterface.(map[string]interface{})
    username := claims["preferredusername"].(string)

    _, err := database.DB.Exec(
        "INSERT INTO click_logs (affiliator_id, item_clicked, referrer_url) VALUES ((SELECT id FROM clients WHERE name = $1), $2, $3)",
        username, clickLog.ItemClicked, clickLog.ReferrerURL,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log click"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Click logged successfully"})
}
