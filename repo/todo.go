package repo

import (
	"todotdd/models"

	"gorm.io/gorm"
)

func UpdateTodo(db *gorm.DB, todoToSave models.Todo) error {
	if err := db.First(&models.Todo{}, todoToSave.ID).Updates(todoToSave).Error; err != nil {
		return models.TodoErr{Code: models.DbErr, Message: err.Error()}
	}
	return nil
}
