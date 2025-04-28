package main

import (
    "myproject/auth"
    "myproject/database"
    "myproject/handlers"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

func main() {
    database.ConnectDB()
    auth.InitCasbin()

    r := gin.Default()

    //  เพิ่ม CORS ให้ React (localhost:3000) เรียกได้
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

    //  Routes
    authGroup := r.Group("/api")
    authGroup.Use(auth.JWTAuthMiddleware(), handlers.LogMiddleware())
    {
        authGroup.GET("/data", handlers.GetData)
    }

    // สมัคร Affiliator Website (optional ทำต่อได้)
    // authGroup.POST("/affiliator/register", handlers.RegisterWebsite)

    r.Run(":8080")
}