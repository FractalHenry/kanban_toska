package repository

import (
	"backend/models"
	"fmt"
)

// Функция создания задания
func (r *Repository) CreateTask(task *models.Task, userLogin string) error {
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Create(task).Error
	}
	return fmt.Errorf("у пользователя нет прав для создания задания на этой доске")
}

// Функция обновления задания
func (r *Repository) UpdateTask(task *models.Task, userLogin string) error {
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Save(task).Error
	}
	return fmt.Errorf("у пользователя нет прав для обновления задания на этой доске")
}

// Функция архивирования задания
func (r *Repository) ArchiveTask(taskID uint, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		task.TaskInArchive = true
		return r.db.Save(&task).Error
	}
	return fmt.Errorf("у пользователя нет прав для архивирования задания")
}

// Функция разархивирования задания
func (r *Repository) UnarchiveTask(taskID uint, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		task.TaskInArchive = false
		return r.db.Save(&task).Error
	}
	return fmt.Errorf("у пользователя нет прав для разархивации задания")
}

// Функция удаления задания
func (r *Repository) DeleteTask(taskID uint, userLogin string) error {
	task, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Delete(&task).Error
	}
	return fmt.Errorf("у пользователя нет прав для удаления задания")
}

// Функция получения всех заданий для карты
func (r *Repository) GetCardTasks(cardID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Where("card_id = ?", cardID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// Вспомогательная функция для получения задания по ID
func (r *Repository) getTaskByID(taskID uint) (models.Task, error) {
	var task models.Task
	if err := r.db.First(&task, taskID).Error; err != nil {
		return task, err
	}
	return task, nil
}
