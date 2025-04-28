package handlers

import (
    "myproject/database"
    "myproject/models"
    "net/http"
    "github.com/gin-gonic/gin"
)

func GetHotels(c *gin.Context) {
    location := c.Query("location")

    query := "SELECT id, name, location, price_per_night, thumbnail_url, detailurl FROM hotels"
    var args []interface{}

    if location != "" {
        query += " WHERE LOWER(location) LIKE LOWER($1)"
        args = append(args, "%"+location+"%")
    }

    rows, err := database.DB.Query(query, args...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch hotels"})
        return
    }
    defer rows.Close()

    hotels := []models.Hotel{}

    for rows.Next() {
        var hotel models.Hotel
        if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Location, &hotel.PricePerNight, &hotel.ThumbnailURL, &hotel.DetailURL); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning hotel"})
            return
        }
        hotels = append(hotels, hotel)
    }

    claimsInterface, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    claims := claimsInterface.(map[string]interface{})
    username := claims["preferredusername"].(string)

    _, err = database.DB.Exec(
        "INSERT INTO request_logs (affiliator_id, endpoint, method, parameters) VALUES ((SELECT id FROM clients WHERE name = $1), $2, $3, $4)",
        username, "/api/hotels", "GET", "location="+location,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log request"})
        return
    }

    c.JSON(http.StatusOK, hotels)
}
