package handlers

import (
	"net/http"
	"strconv"
	"todoapp/models"

	"github.com/gin-gonic/gin"
)

type TodoDto struct {
	Id          int    `json:"id" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsCompleted bool   `json:"is_completed" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
}

func UpdateTodo(c *gin.Context) {
	var todoDto TodoDto
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.TodoErr{Code: models.IdNotIntegerErr, Message: err.Error()})
		return
	}
	todoDto.Id = id
	if err := c.ShouldBind(&todoDto); err != nil {
		c.JSON(http.StatusBadRequest, models.TodoErr{Code: models.ValidationErr, Message: err.Error()})
		return
	}
	c.JSON(http.StatusNotFound, nil)
}
