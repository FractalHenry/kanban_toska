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
	task, err := r.GetTaskByID(taskID)
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
	task, err := r.GetTaskByID(taskID)
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
	task, err := r.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		// Удаляем все связанные данные, игнорируя ошибки отсутствия записей
		_ = r.db.Where("task_id = ?", taskID).Delete(&models.Mark{}).Error
		_ = r.db.Where("task_id = ?", taskID).Delete(&models.MarkName{}).Error
		_ = r.db.Where("task_id = ?", taskID).Delete(&models.Checklist{}).Error
		_ = r.db.Where("checklist_id IN (?)", r.db.Model(&models.Checklist{}).Where("task_id = ?", taskID).Select("checklist_id")).Delete(&models.ChecklistElement{}).Error
		_ = r.db.Where("task_id = ?", taskID).Delete(&models.TaskColor{}).Error
		_ = r.db.Where("task_id = ?", taskID).Delete(&models.TaskDescription{}).Error
		_ = r.db.Where("task_id = ?", taskID).Delete(&models.TaskDateStart{}).Error
		_ = r.db.Where("task_id = ?", taskID).Delete(&models.TaskDateEnd{}).Error
		_ = r.db.Where("task_date_end_id IN (?)", r.db.Model(&models.TaskDateEnd{}).Where("task_id = ?", taskID).Select("task_date_end_id")).Delete(&models.TaskNotification{}).Error

		// Удаляем само задание
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
func (r *Repository) GetTaskByID(taskID uint) (models.Task, error) {
	var task models.Task
	if err := r.db.First(&task, taskID).Error; err != nil {
		return task, err
	}
	return task, nil
}

func (r *Repository) GetTaskNotifications(taskID uint) (*[]models.Notification, error) {
	var notifications []models.Notification
	var taskDateEnd models.TaskDateEnd

	err := r.db.Where("task_id = ?", taskID).First(&taskDateEnd).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&models.TaskNotification{}).
		Where("task_date_end_id = ?", taskDateEnd.TaskDateEndID).
		Association("Notification").
		Find(&notifications)

	if err != nil {
		return nil, err
	}

	return &notifications, nil
}

func (r *Repository) GetTaskMarks(taskID uint) (*[]models.Mark, error) {
	var marks []models.Mark
	err := r.db.Where("task_id = ?", taskID).Find(&marks).Error
	if err != nil {
		return nil, err
	}
	return &marks, nil
}

func (r *Repository) GetTaskMarkNames(taskID uint) (*[]models.MarkName, error) {
	var markNames []models.MarkName
	err := r.db.Model(&models.Mark{}).
		Where("task_id = ?", taskID).
		Association("MarkName").
		Find(&markNames)

	if err != nil {
		return nil, err
	}

	return &markNames, nil
}
