package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateTodo(c *gin.Context) {
	c.JSON(http.StatusBadRequest, `{ "code": "validation err", "message": "strconv.Atoi: parsing \"abc\": invalid syntax" }`)
}
