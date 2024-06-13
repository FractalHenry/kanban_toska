package repository

import (
	"backend/models"
	"fmt"
)

// Функция создания чек-листа
func (r *Repository) CreateChecklist(checklist *models.Checklist, userLogin string) error {
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(checklist, userLogin)
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
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(checklist, userLogin)
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

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&checklist, userLogin)
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
		return &checklist, err
	}
	return &checklist, nil
}

// Функция создания элемента чек-листа
func (r *Repository) CreateChecklistElement(element *models.ChecklistElement, userLogin string) error {
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(element, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Create(element).Error
	}
	return fmt.Errorf("у пользователя нет прав для создания элемента чек-листа на этой доске")
}

// Функция обновления элемента чек-листа
func (r *Repository) UpdateChecklistElement(element *models.ChecklistElement, userLogin string) error {
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(element, userLogin)
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

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&element, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Delete(&element).Error
	}
	return fmt.Errorf("у пользователя нет прав для удаления элемента чек-листа")
}

// Вспомогательная функция для получения элемента чек-листа по ID
func (r *Repository) GetChecklistElementByID(elementID uint) (models.ChecklistElement, error) {
	var element models.ChecklistElement
	if err := r.db.First(&element, elementID).Error; err != nil {
		return element, err
	}
	return element, nil
}
