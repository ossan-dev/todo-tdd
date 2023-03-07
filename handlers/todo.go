package handlers

import (
	"net/http"
	"todoapp/models"

	"github.com/gin-gonic/gin"
)

type TodoErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func UpdateTodo(c *gin.Context) {
	c.JSON(http.StatusBadRequest, TodoErr{Code: models.IdNotIntegerErr, Message: "strconv.Atoi: parsing \"abc\": invalid syntax\" }"})
}
