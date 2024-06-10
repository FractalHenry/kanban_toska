package repository

import (
	"backend/models"
	"fmt"
)

// Функция создания карточки
func (r *Repository) CreateCard(card *models.Card, userLogin string) error {
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(card, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Create(card).Error
	}
	return fmt.Errorf("у пользователя нет прав для создания карточки на этой доске")
}

// Функция обновления карточки
func (r *Repository) UpdateCard(card *models.Card, userLogin string) error {
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(card, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Save(card).Error
	}
	return fmt.Errorf("у пользователя нет прав для обновления карточки на этой доске")
}

// Функция архивирования карточки
func (r *Repository) ArchiveCard(cardID uint, userLogin string) error {
	card, err := r.getCardByID(cardID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&card, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		card.CardInArchive = true
		return r.db.Save(&card).Error
	}
	return fmt.Errorf("у пользователя нет прав для архивирования карточки")
}

// Функция разархивирования карточки
func (r *Repository) UnarchiveCard(cardID uint, userLogin string) error {
	card, err := r.getCardByID(cardID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&card, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		card.CardInArchive = false
		return r.db.Save(&card).Error
	}
	return fmt.Errorf("у пользователя нет прав для разархивирования карточки")
}

// Функция удаления карточки
func (r *Repository) DeleteCard(cardID uint, userLogin string) error {
	card, err := r.getCardByID(cardID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&card, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Delete(&card).Error
	}
	return fmt.Errorf("у пользователя нет прав для удаления карточки")
}

// Функция получения всех карточек для доски
func (r *Repository) GetBoardCards(boardID uint) ([]models.Card, error) {
	var cards []models.Card
	if err := r.db.Where("board_id = ?", boardID).Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

// Вспомогательная функция для получения карточки по ID
func (r *Repository) getCardByID(cardID uint) (models.Card, error) {
	var card models.Card
	if err := r.db.First(&card, cardID).Error; err != nil {
		return card, err
	}
	return card, nil
}
