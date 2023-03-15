package handlers

import (
	"net/http"
	"strconv"

	"todotdd/models"
	"todotdd/repo"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const DBKey = "DB"

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

	todo := models.NewTodo(todoDto.Id, todoDto.Description, todoDto.IsCompleted, todoDto.DueDate)

	db := c.MustGet(DBKey).(*gorm.DB)
	if db != nil {
		err = repo.UpdateTodo(db, todo)
		if err != nil {
			todoErr := err.(models.TodoErr)
			c.JSON(todoErr.StatusCode, todoErr)
			return
		}
	}

	c.JSON(http.StatusAccepted, nil)
}
