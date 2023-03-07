package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todotdd/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestUpdateTodo_IdNotInteger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "abcd"})

	UpdateTodo(c)

	var todoErr models.TodoErr
	if err := json.Unmarshal(w.Body.Bytes(), &todoErr); err != nil {
		t.Fatalf("err not expected while unmarshaling: %v", err)
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, models.IdNotIntegerErr, todoErr.Code)
}

func TestUpdateTodo_ValidationErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPut, "/todos", strings.NewReader(`{ "description": "", "is_completed": true, "due_date": "2023-05-04" }`))
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request = req

	UpdateTodo(c)

	var todoErr models.TodoErr
	if err := json.Unmarshal(w.Body.Bytes(), &todoErr); err != nil {
		t.Fatalf("err not expected while unmarshaling: %v", err)
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, models.ValidationErr, todoErr.Code)
}

func TestUpdateTodo_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPut, "/todos", strings.NewReader(`{ "description": "lorem ipsum", "is_completed": true, "due_date": "2023-05-04" }`))
	req.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request = req

	UpdateTodo(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateTodo_NotFound_Body(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPut, "/todos", strings.NewReader(`{ "description": "lorem ipsum", "is_completed": true, "due_date": "2023-05-04" }`))
	req.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request = req

	UpdateTodo(c)

	var todoErr models.TodoErr
	if err := json.Unmarshal(w.Body.Bytes(), &todoErr); err != nil {
		t.Fatalf("err not expected while unmarshaling: %v", err)
	}
	assert.Equal(t, "unknown todo", todoErr.Code)
}
