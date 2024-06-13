package repository

import (
	"backend/models"
	"fmt"
	"time"
)

// Функция установки даты начала задачи
func (r *Repository) SetTaskStartDate(taskID uint, startDate time.Time, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		taskDateStart := models.TaskDateStart{
			TaskID:        taskID,
			TaskDateStart: startDate,
		}
		// Проверяем существует ли запись для taskID
		if err := r.db.Where("task_id = ?", taskID).First(&taskDateStart).Error; err != nil {
			// Если записи не существует, создаем новую
			return r.db.Create(&taskDateStart).Error
		}
		// Обновляем существующую запись
		taskDateStart.TaskDateStart = startDate
		return r.db.Save(&taskDateStart).Error
	}
	return fmt.Errorf("у пользователя нет прав для установки даты начала задачи")
}

// Функция удаления даты начала задачи
func (r *Repository) RemoveTaskStartDate(taskID uint, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		if err := r.db.Where("task_id = ?", taskID).Delete(&models.TaskDateStart{}).Error; err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("у пользователя нет прав для удаления даты начала задачи")
}

// Функция установки даты окончания задачи
func (r *Repository) SetTaskEndDate(taskID uint, endDate time.Time, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		taskDateEnd := models.TaskDateEnd{
			TaskID:      taskID,
			TaskDateEnd: endDate,
		}
		// Проверяем существует ли запись для taskID
		if err := r.db.Where("task_id = ?", taskID).First(&taskDateEnd).Error; err != nil {
			// Если записи не существует, создаем новую
			return r.db.Create(&taskDateEnd).Error
		}
		// Обновляем существующую запись
		taskDateEnd.TaskDateEnd = endDate
		return r.db.Save(&taskDateEnd).Error
	}
	return fmt.Errorf("у пользователя нет прав для установки даты окончания задачи")
}

// Функция удаления даты окончания задачи
func (r *Repository) RemoveTaskEndDate(taskID uint, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		if err := r.db.Where("task_id = ?", taskID).Delete(&models.TaskDateEnd{}).Error; err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("у пользователя нет прав для удаления даты окончания задачи")
}

// Функция получения даты начала задачи
func (r *Repository) GetTaskStartDate(taskID uint) (*models.TaskDateStart, error) {
	var taskDateStart models.TaskDateStart
	if err := r.db.Where("task_id = ?", taskID).First(&taskDateStart).Error; err != nil {
		return nil, err
	}
	return &taskDateStart, nil
}

// Функция получения даты окончания задачи
func (r *Repository) GetTaskEndDate(taskID uint) (*models.TaskDateEnd, error) {
	var taskDateEnd models.TaskDateEnd
	if err := r.db.Where("task_id = ?", taskID).First(&taskDateEnd).Error; err != nil {
		return nil, err
	}
	return &taskDateEnd, nil
}
