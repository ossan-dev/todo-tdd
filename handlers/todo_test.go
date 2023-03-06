package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoapp/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestUpdateTodo_IdNotInteger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "abcd"})

	UpdateTodo(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateTodo_IdNotInteger_Body(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var todoErr TodoErr

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "abcd"})

	UpdateTodo(c)

	if err := json.Unmarshal(w.Body.Bytes(), &todoErr); err != nil {
		t.Fatalf("err not expected while unmarshaling: %v", err)
	}

	assert.Equal(t, models.ValidationErr, todoErr.Code)
}
