package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	expectedBody := `"{ \"code\": \"validation err\", \"message\": \"strconv.Atoi: parsing \\\"abc\\\": invalid syntax\" }"`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "abcd"})

	UpdateTodo(c)

	if w.Body.String() != expectedBody {
		t.Errorf("expected %q got %q", expectedBody, w.Body.String())
	}
}
