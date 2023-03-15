package repo

import (
	"errors"
	"net/http"

	"todotdd/models"

	"gorm.io/gorm"
)

func UpdateTodo(db *gorm.DB, todoToSave models.Todo) error {
	if err := db.First(&models.Todo{}, todoToSave.ID).Updates(todoToSave).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TodoErr{StatusCode: http.StatusNotFound, Code: models.TodoNotFoundErr, Message: err.Error()}
		}
		return models.TodoErr{StatusCode: http.StatusInternalServerError, Code: models.DbErr, Message: err.Error()}
	}
	return nil
}
