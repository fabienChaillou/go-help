package main

import (
	// "yourapp/config"
	// "yourapp/handler"
	// "yourapp/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db := InitDB()
	r := gin.Default()

	r.GET("/users", GetUsers(db))
	r.GET("/users/:id", GetUserByID(db))

	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	auth.POST("/users", CreateUser(db))
	auth.PUT("/users/:id", UpdateUser(db))
	auth.DELETE("/users/:id", DeleteUser(db))

	r.Run(":8080")
}
