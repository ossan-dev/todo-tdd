package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUpdateTodo_IdNotInteger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	UpdateTodo(c)

	if w.Code != 400 {
		t.Fatalf("expected %d got %d", 400, w.Code)
	}
}
