package handlers

import (
	"net/http"
	"strconv"

	"todotdd/models"
	"todotdd/repo"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const DbKey = "Db"

type TodoDto struct {
	Description string `json:"description" binding:"required"`
	IsCompleted bool   `json:"is_completed" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
}

func UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.TodoErr{Code: models.IdNotIntegerErr, Message: err.Error()})
		return
	}
	var todoDto TodoDto
	if err := c.ShouldBind(&todoDto); err != nil {
		c.JSON(http.StatusBadRequest, models.TodoErr{Code: models.ValidationErr, Message: err.Error()})
		return
	}
	todoToSave := models.Todo{Description: todoDto.Description, IsCompleted: todoDto.IsCompleted, DueDate: todoDto.DueDate}
	gormDb := c.MustGet(DbKey).(*gorm.DB)
	if err = repo.UpdateTodo(gormDb, id, todoToSave); err != nil {
		todoErr := err.(models.TodoErr)
		c.JSON(todoErr.StatusCode, todoErr)
		return
	}
	c.Writer.WriteHeader(http.StatusAccepted)
	c.Writer.WriteHeaderNow()
}
