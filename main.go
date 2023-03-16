package main

import (
	"todotdd/handlers"
	"todotdd/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=54322 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Todo{})

	// db bootstrap
	todos := []models.Todo{
		{ID: 1, Description: "first todo", IsCompleted: false, DueDate: "2023-06-01"},
		{ID: 2, Description: "second todo", IsCompleted: true, DueDate: "2023-01-15"},
		{ID: 3, Description: "third todo", IsCompleted: false, DueDate: "2023-02-28"},
	}
	db.Create(todos)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	})

	r.PUT("todos/:id", handlers.UpdateTodo)

	r.Run()
}
