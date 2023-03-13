package handlers

import (
	"net/http"
	"strconv"

	"todotdd/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TodoDto struct {
	Id          int    `json:"id" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsCompleted bool   `json:"is_completed" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
}

type Todo struct {
	ID          int
	Description string
	IsCompleted bool
	DueDate     string
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

	todo := Todo{ID: todoDto.Id, Description: todoDto.Description, IsCompleted: todoDto.IsCompleted, DueDate: todoDto.DueDate}

	dbKey, isFound := c.Keys["DB"]
	if isFound {
		db := dbKey.(*gorm.DB)
		if db != nil {
			err = db.First(&Todo{}, todo.ID).Updates(todo).Error
			if err != nil {
				c.JSON(500, err)
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, models.TodoErr{Code: models.TodoNotFoundErr, Message: "unknown todo"})
}
