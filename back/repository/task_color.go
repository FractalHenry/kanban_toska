package repository

import (
	"backend/models"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

func (r *Repository) CreateTaskColor(taskColor *models.TaskColor, userLogin string) error {
	task, err := r.getTaskByID(taskColor.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Save(&taskColor).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для создания цвета задания")
	}
}

func (r *Repository) DeleteTaskColor(taskID uint, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		taskColor := models.TaskColor{
			TaskID: taskID,
		}
		return r.db.Delete(&taskColor).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления цвета задания")
	}
}

func (r *Repository) GetTaskColor(taskID uint) (string, error) {
	var taskColor models.TaskColor
	err := r.db.Where("task_id = ?", taskID).First(&taskColor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // Если цвета не существует, возвращаем пустую строку
		}
		return "", err
	}
	return taskColor.TaskColor, nil
}
