package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateTodo(c *gin.Context) {
	c.JSON(http.StatusBadRequest, nil)
}
