package main

import (
	"context"
	"time"
	"todotdd/handlers"
	"todotdd/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost port=54322 user=postgres password=postgres dbname=postgres sslmode=disable"
	gormDb, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	gormDb.AutoMigrate(&models.Todo{})

	gormDb.Create(&models.Todo{Id: 1, Description: "sample todo", IsCompleted: false, DueDate: "2023-02-15"})

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		timeoutCtx, cancelFunc := context.WithTimeout(ctx.Request.Context(), time.Second*5)
		defer cancelFunc()

		ctx.Set(handlers.DbKey, gormDb.WithContext(timeoutCtx))
		ctx.Next()
	})

	r.PUT("/todos/:id", handlers.UpdateTodo)

	r.Run(":8000")
}
