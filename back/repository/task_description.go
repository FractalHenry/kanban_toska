package repository

import (
	"backend/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateTaskDescription(taskDescription *models.TaskDescription, userLogin string) error {
	task, err := r.GetTaskByID(taskDescription.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Create(taskDescription).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для создания описания задания")
	}
}

func (r *Repository) UpdateTaskDescription(taskDescription *models.TaskDescription, userLogin string) error {
	task, err := r.GetTaskByID(taskDescription.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "task_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"task_description"}),
		}).Create(taskDescription).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для изменения описания задания")
	}
}

func (r *Repository) DeleteTaskDescription(taskID uint, userLogin string) error {
	task, err := r.GetTaskByID(taskID)
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

		result := r.db.Unscoped().Where("task_id = ?", taskID).Delete(&taskDescription)
		if result.RowsAffected == 0 {
			// Запись не найдена, но это не ошибка
			return nil
		}

		if result.Error != nil {
			return result.Error
		}

		return nil
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
