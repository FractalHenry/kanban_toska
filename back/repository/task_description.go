package repository

import (
	"backend/models"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

func (r *Repository) UpdateTaskDescription(taskID uint, newDescription string, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		taskDescription := models.TaskDescription{
			TaskDescription: newDescription,
			TaskID:          taskID,
		}
		return r.db.Save(&taskDescription).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для изменения описания задания")
	}
}

func (r *Repository) DeleteTaskDescription(taskID uint, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		taskDescription := models.TaskDescription{
			TaskID: taskID,
		}
		return r.db.Delete(&taskDescription).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления описания задания")
	}
}

func (r *Repository) GetTaskDescription(taskID uint) (string, error) {
	var taskDescription models.TaskDescription
	err := r.db.Where("task_id = ?", taskID).First(&taskDescription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // Если описания не существует, возвращаем пустую строку
		}
		return "", err
	}
	return taskDescription.TaskDescription, nil
}
