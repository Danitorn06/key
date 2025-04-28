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
    // เชื่อมต่อ Database
    database.ConnectDB()

    // โหลด Casbin (สำหรับเช็กสิทธิ์ role)
    auth.InitCasbin()

    // สร้าง Gin Engine
    r := gin.Default()

    // ตั้งค่า CORS เพื่อให้ Frontend ที่ localhost:3000 เรียกได้
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

    // ===== Public Routes =====
    r.POST("/register", handlers.RegisterClient) // ยังเปิดไว้ สำหรับ Affiliator ลงทะเบียนในระบบ

    // ===== Protected Routes (ต้องมี Keycloak Token) =====
    authGroup := r.Group("/api")
    authGroup.Use(auth.JWTAuthMiddleware(), handlers.LogMiddleware())
    {
        // Affiliator Website
        authGroup.POST("/affiliator/register-website", handlers.RegisterWebsite)

        // Hotels
        authGroup.GET("/hotels", handlers.GetHotels)

        // Click Logs
        authGroup.POST("/click-log", handlers.LogClick)

        // Request Logs
        authGroup.GET("/request-logs", handlers.GetRequestLogs)
    }

    // ===== Start Server =====
    r.Run(":8080")
}