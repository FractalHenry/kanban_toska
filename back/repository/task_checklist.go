package repository

import (
	"backend/models"
	"fmt"
)

// Функция для получения всех элементов чек-листа по ID чек-листа
func (r *Repository) GetChecklistElementsByChecklistID(checklistID uint) (*[]models.ChecklistElement, error) {
	var elements *[]models.ChecklistElement
	if err := r.db.Where("checklist_id = ?", checklistID).Find(&elements).Error; err != nil {
		return nil, err
	}
	return elements, nil
}

// Функция для получения всех чек-листов по ID задачи
func (r *Repository) GetChecklistsByTaskID(taskID uint) (*[]models.Checklist, error) {
	var checklists *[]models.Checklist
	if err := r.db.Where("task_id = ?", taskID).Find(&checklists).Error; err != nil {
		return nil, err
	}
	return checklists, nil
}

// Функция создания чек-листа
func (r *Repository) CreateChecklist(checklist *models.Checklist, userLogin string) error {
	task, err := r.GetTaskByID(checklist.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Create(checklist).Error
	}
	return fmt.Errorf("у пользователя нет прав для создания чек-листа на этой доске")
}

// Функция обновления чек-листа
func (r *Repository) UpdateChecklist(checklist *models.Checklist, userLogin string) error {
	task, err := r.GetTaskByID(checklist.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Save(checklist).Error
	}
	return fmt.Errorf("у пользователя нет прав для обновления чек-листа на этой доске")
}

// Функция удаления чек-листа
func (r *Repository) DeleteChecklist(checklistID uint, userLogin string) error {
	checklist, err := r.GetChecklistByID(checklistID)
	if err != nil {
		return err
	}

	task, err := r.GetTaskByID(checklist.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Delete(&checklist).Error
	}
	return fmt.Errorf("у пользователя нет прав для удаления чек-листа")
}

// Вспомогательная функция для получения чек-листа по ID
func (r *Repository) GetChecklistByID(checklistID uint) (*models.Checklist, error) {
	var checklist models.Checklist
	if err := r.db.First(&checklist, checklistID).Error; err != nil {
		return nil, err
	}
	return &checklist, nil
}

// Функция создания элемента чек-листа
func (r *Repository) CreateChecklistElement(element *models.ChecklistElement, userLogin string) error {
	checklist, err := r.GetChecklistByID(element.ChecklistID)
	if err != nil {
		return err
	}

	task, err := r.GetTaskByID(checklist.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		// Получаем следующий порядковый номер для нового элемента
		order, err := r.getNextChecklistElementOrder(checklist.ChecklistID)
		if err != nil {
			return err
		}
		element.Order = order

		return r.db.Create(element).Error
	}
	return fmt.Errorf("у пользователя нет прав для создания элемента чек-листа на этой доске")
}

// Функция обновления элемента чек-листа
func (r *Repository) UpdateChecklistElement(element *models.ChecklistElement, userLogin string) error {
	checklist, err := r.GetChecklistByID(element.ChecklistID)
	if err != nil {
		return err
	}

	task, err := r.GetTaskByID(checklist.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Save(element).Error
	}
	return fmt.Errorf("у пользователя нет прав для обновления элемента чек-листа на этой доске")
}

// Функция удаления элемента чек-листа
func (r *Repository) DeleteChecklistElement(elementID uint, userLogin string) error {
	element, err := r.GetChecklistElementByID(elementID)
	if err != nil {
		return err
	}

	checklist, err := r.GetChecklistByID(element.ChecklistID)
	if err != nil {
		return err
	}

	task, err := r.GetTaskByID(checklist.TaskID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&task, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		// Обновляем порядковые номера остальных элементов чеклиста
		err = r.updateChecklistElementOrders(element.ChecklistID, element.Order)
		if err != nil {
			return err
		}

		return r.db.Delete(&element).Error
	}
	return fmt.Errorf("у пользователя нет прав для удаления элемента чек-листа")
}

// Вспомогательная функция для получения элемента чек-листа по ID
func (r *Repository) GetChecklistElementByID(elementID uint) (*models.ChecklistElement, error) {
	var element models.ChecklistElement
	if err := r.db.First(&element, elementID).Error; err != nil {
		return nil, err
	}
	return &element, nil
}

// Вспомогательная функция для получения следующего порядкового номера для элемента чек-листа
func (r *Repository) getNextChecklistElementOrder(checklistID uint) (uint, error) {
	var maxOrder uint
	err := r.db.Model(&models.ChecklistElement{}).
		Where("checklist_id = ?", checklistID).
		Select("COALESCE(MAX(order), 0)").
		Scan(&maxOrder).Error
	if err != nil {
		return 0, err
	}
	return maxOrder + 1, nil
}

// Вспомогательная функция для обновления порядковых номеров элементов чек-листа после удаления одного из них
func (r *Repository) updateChecklistElementOrders(checklistID, deletedOrder uint) error {
	var elements []models.ChecklistElement
	err := r.db.Where("checklist_id = ? AND order > ?", checklistID, deletedOrder).
		Order("order ASC").
		Find(&elements).Error
	if err != nil {
		return err
	}

	for _, element := range elements {
		element.Order--
		err = r.db.Save(&element).Error
		if err != nil {
			return err
		}
	}

	return nil
}
