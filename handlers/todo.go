package handlers

import (
	"net/http"
	"strconv"
	"todoapp/models"

	"github.com/gin-gonic/gin"
)

type TodoErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func UpdateTodo(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, TodoErr{Code: models.IdNotIntegerErr, Message: "strconv.Atoi: parsing \"abc\": invalid syntax\" }"})
		return
	}
	c.JSON(http.StatusBadRequest, TodoErr{Code: "validation err", Message: "strconv.Atoi: parsing \"abc\": invalid syntax\" }"})
}
