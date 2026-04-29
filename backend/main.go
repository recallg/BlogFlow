package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/articles", getArticleListHandler(DB))
		api.GET("/articles/:id", getArticleHandler(DB))
		api.POST("/articles", createArticleHandler(DB))
		api.PUT("/articles/:id", updateArticleHandler(DB))
		api.DELETE("/articles/:id", deleteArticleHandler(DB))

		// Auth routes
		api.POST("/auth/register", registerHandler(DB))
		api.POST("/auth/login", loginHandler(DB))
		api.GET("/auth/me", AuthMiddleware(), getMeHandler(DB))
		api.PUT("/auth/change-password", AuthMiddleware(), changePasswordHandler(DB))
	}

	log.Println("🚀 Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
