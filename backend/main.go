package main

import (
	"backend/controllers"
	"backend/db"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	router := gin.Default()

	router.Use(middlewares.CorsMiddleware())

	//Connect to db and send schema
	db.ConnectDb()
	db.MigrateDb()

	//Route for register
	router.POST("/register", controllers.Register())
	router.POST("/login", controllers.Login())

	auth := router.Group("/").Use(middlewares.Authentication())

	{
		auth.POST("/users", middlewares.CreateUser())
		auth.POST("/logout", controllers.Logout())
	}

	router.Run(":8000")
}
