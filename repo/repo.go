package repo

import (
	"errors"
	"net/http"
	"todotdd/models"

	"gorm.io/gorm"
)

func UpdateTodo(db *gorm.DB, id int, todoToSave models.Todo) error {
	// these two lines are used to emulate slowness on the DB side
	// they give you the context deadline exceeded error
	// var test string
	// db.Raw("SELECT pg_sleep(10) as test").Scan(&test)
	if err := db.Model(&models.Todo{}).First(&models.Todo{}, id).Updates(&todoToSave).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TodoErr{StatusCode: http.StatusNotFound, Code: models.NotFoundErr, Message: err.Error()}
		}
		return models.TodoErr{StatusCode: http.StatusInternalServerError, Code: models.DbErr, Message: err.Error()}
	}
	return nil
}
