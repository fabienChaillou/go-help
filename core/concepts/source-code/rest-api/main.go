// Structure de projet Go modulaire (Gin + PostgreSQL + Generics + Auth + Validation)

// --- main.go ---
package main

import (
	"rest-api/config"
	"rest-api/handler"
	"rest-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	r := gin.Default()

	r.GET("/users", handler.GetUsers(db))
	r.GET("/users/:id", handler.GetUserByID(db))

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/users", handler.CreateUser(db))
	auth.PUT("/users/:id", handler.UpdateUser(db))
	auth.DELETE("/users/:id", handler.DeleteUser(db))

	r.Run(":8080")
}
