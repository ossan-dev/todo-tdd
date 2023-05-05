package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todotdd/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 1. id not integer
// 2. validation error
// 3. db error
// 4. todo not found
// 5. happy path

func TestUpdateTodo_IdNotInteger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c := gin.CreateTestContextOnly(w, gin.Default())

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "abc"})

	UpdateTodo(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), models.IdNotIntegerErr)
}

func TestUpdateTodo_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c := gin.CreateTestContextOnly(w, gin.Default())

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Method = http.MethodPut
	c.Request.Body = io.NopCloser(strings.NewReader(`{"description": "", "is_completed": true, "due_date": "2023-05-05"}`))

	UpdateTodo(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), models.ValidationErr)
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
	}))

	mock.ExpectQuery(`SELECT * FROM "todos" WHERE "todos"."id" = $1 ORDER BY "todos"."id" LIMIT 1`).WithArgs(1).WillReturnError(fmt.Errorf("db error"))

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Method = http.MethodPut
	c.Request.Body = io.NopCloser(strings.NewReader(`{"description": "lorem ipsum", "is_completed": true, "due_date": "2023-05-05"}`))
	c.Set(DbKey, gormDb)

	UpdateTodo(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), models.DbErr)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}
}

func TestUpdateTodo_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c := gin.CreateTestContextOnly(w, gin.Default())
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()

	conn, _ := db.Conn(c)
	defer conn.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}))

	mock.ExpectQuery(`SELECT * FROM "todos" WHERE "todos"."id" = $1 ORDER BY "todos"."id" LIMIT 1`).WithArgs(1).WillReturnError(gorm.ErrRecordNotFound)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Method = http.MethodPut
	c.Request.Body = io.NopCloser(strings.NewReader(`{"description": "lorem ipsum", "is_completed": true, "due_date": "2023-05-05"}`))
	c.Set(DbKey, gormDb)

	UpdateTodo(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), models.NotFoundErr)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
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
	}))

	rows := sqlmock.NewRows([]string{"id", "description", "is_completed", "due_date"}).AddRow(1, "sample todo", false, "2023-02-15")

	mock.ExpectQuery(`SELECT * FROM "todos" WHERE "todos"."id" = $1 ORDER BY "todos"."id" LIMIT 1`).WithArgs(1).WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "todos" SET "description"=$1,"is_completed"=$2,"due_date"=$3 WHERE "todos"."id" = $4`).WithArgs("lorem ipsum", true, "2023-05-05", 1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Method = http.MethodPut
	c.Request.Body = io.NopCloser(strings.NewReader(`{"description": "lorem ipsum", "is_completed": true, "due_date": "2023-05-05"}`))
	c.Set(DbKey, gormDb)

	UpdateTodo(c)

	assert.Equal(t, http.StatusAccepted, w.Code)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}
}
