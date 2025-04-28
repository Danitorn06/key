package handlers

import (
    "myproject/database"
    "myproject/models"
    "net/http"

    "github.com/gin-gonic/gin"
)

// GetRequestLogs ดึงประวัติการ Request ของ Affiliator
func GetRequestLogs(c *gin.Context) {
    claims, _ := c.Get("claims")
    mapClaims := claims.(map[string]interface{})
    username := mapClaims["preferred_username"].(string)

    // หา affiliator_id
    var affiliatorID int
    err := database.DB.QueryRow("SELECT id FROM clients WHERE name = $1", username).Scan(&affiliatorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Affiliator not found"})
        return
    }

    // ดึง request logs
    rows, err := database.DB.Query(
        "SELECT id, endpoint, method, parameters, requested_at FROM request_logs WHERE affiliator_id = $1 ORDER BY requested_at DESC",
        affiliatorID,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch request logs"})
        return
    }
    defer rows.Close()

    var logs []models.RequestLog

    for rows.Next() {
        var log models.RequestLog
        err := rows.Scan(&log.ID, &log.Endpoint, &log.Method, &log.Parameters, &log.RequestedAt)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning logs"})
            return
        }
        logs = append(logs, log)
    }

    c.JSON(http.StatusOK, logs)
}