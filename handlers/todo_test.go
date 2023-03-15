package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"todotdd/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	var todoErr models.TodoErr
	if err := json.Unmarshal(w.Body.Bytes(), &todoErr); err != nil {
		t.Fatalf("err not expected while unmarshaling: %v", err)
	}
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, models.TodoNotFoundErr, todoErr.Code)
}

func TestUpdateTodo_DbError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c := gin.CreateTestContextOnly(w, gin.Default())

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	conn, _ := db.Conn(c)
	defer conn.Close()
	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{})

	mock.ExpectQuery(`SELECT * FROM "todos" WHERE "todos"."id" = $1 ORDER BY "todos"."id" LIMIT 1`).WithArgs(1).WillReturnError(fmt.Errorf("db error"))

	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	c.Request.Method = http.MethodPut
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{ "description": "lorem ipsum", "is_completed": true, "due_date": "2023-05-04" }`))
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Set(DBKey, gormDb)

	UpdateTodo(c)

	var todoErr models.TodoErr
	if err := json.Unmarshal(w.Body.Bytes(), &todoErr); err != nil {
		t.Fatalf("err not expected while unmarshaling: %v", err)
	}
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, models.DbErr, todoErr.Code)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestUpdateTodo_HappyPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c := gin.CreateTestContextOnly(w, gin.Default())

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	conn, _ := db.Conn(c)
	defer conn.Close()
	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{})

	rows := sqlmock.NewRows([]string{"id", "description", "is_completed", "due_date"}).AddRow(1, "lorem ipsum", false, "2023-05-01")
	mock.ExpectQuery(`SELECT * FROM "todos" WHERE "todos"."id" = $1 ORDER BY "todos"."id" LIMIT 1`).WithArgs(1).WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "todos" SET "id"=$1,"description"=$2,"is_completed"=$3,"due_date"=$4 WHERE "todos"."id" = $5 AND "id" = $6`).WithArgs(1, "lorem ipsum", true, "2023-05-04", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	c.Request.Method = http.MethodPut
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{ "description": "lorem ipsum", "is_completed": true, "due_date": "2023-05-04" }`))
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Set(DBKey, gormDb)

	UpdateTodo(c)

	assert.Equal(t, 200, w.Code)
}
