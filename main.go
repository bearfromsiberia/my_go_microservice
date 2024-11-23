package main

import (
	"fmt"
	"log"

	"github.com/bearfromsiberia/my_go_microservice.git/database"
	"github.com/bearfromsiberia/my_go_microservice.git/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	db, err := database.New("./db.sqlite")

	if err != nil {
		fmt.Println("failed to init db")
		fmt.Printf("%s", err)
	}

	r := gin.Default()

	r.POST("/users", handlers.CreateUser(db))
	r.GET("/users", handlers.GetUsers(db))
	r.PATCH("/users/:id", handlers.UpdateUser(db))
	r.DELETE("/users/:id", handlers.DeleteUser(db))
	log.Println("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
