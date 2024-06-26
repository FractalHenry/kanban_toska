package repository

import (
	"backend/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Функция создания карточки
func (r *Repository) CreateCard(card *models.Card, userLogin string) error {
	// Получаем доску, к которой относится карточка
	var board models.Board
	if err := r.db.Where("board_id = ?", card.BoardID).First(&board).Error; err != nil {
		return err
	}

	// Получаем пространство, к которому относится доска
	var space models.Space
	if err := r.db.Where("space_id = ?", board.SpaceID).First(&space).Error; err != nil {
		return err
	}

	// Проверяем роль пользователя на доске
	var userBoardRole models.UserBoardRoleOnBoard
	if err := r.db.Where("login = ? AND board_id = ?", userLogin, card.BoardID).First(&userBoardRole).Error; err == nil {
		// Роль на доске найдена
		var boardRole models.BoardRoleOnBoard
		if err := r.db.Where("role_on_board_id = ?", userBoardRole.RoleOnBoardID).First(&boardRole).Error; err != nil {
			return err
		}

		// Проверяем, есть ли у пользователя права на редактирование доски
		if boardRole.CanEdit {
			return r.db.Create(card).Error
		}

		return fmt.Errorf("у пользователя нет прав для создания карточки на этой доске")
	}

	// Если роль на доске не найдена, проверяем роль пользователя в пространстве
	var roleOnSpace models.RoleOnSpace
	if err := r.db.Table("user_role_on_spaces").
		Select("role_on_spaces.*").
		Joins("JOIN users ON user_role_on_spaces.login = users.login").
		Joins("JOIN role_on_spaces ON user_role_on_spaces.role_on_space_id = role_on_spaces.role_on_space_id").
		Where("users.login = ? AND role_on_spaces.space_id = ?", userLogin, board.SpaceID).
		First(&roleOnSpace).Error; err == nil {

		// Проверяем, есть ли у пользователя права на редактирование пространства
		if roleOnSpace.CanEdit || roleOnSpace.IsAdmin || roleOnSpace.IsOwner {
			return r.db.Create(&card).Error
		}

		return fmt.Errorf("у пользователя нет прав для создания карточки в этом пространстве")
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Запись не найдена, возвращаем ошибку
		return fmt.Errorf("пользователь или пространство не найдены")
	} else {
		// Произошла другая ошибка
		return err
	}
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
	card, err := r.GetCardByID(cardID)
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
	card, err := r.GetCardByID(cardID)
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
	card, err := r.GetCardByID(cardID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRoles(&card, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		// Удаляем все связанные данные
		tasks, err := r.GetCardTasks(cardID)
		if err != nil {
			return err
		}

		for _, task := range tasks {
			r.DeleteTask(task.TaskID, userLogin)
		}

		// Удаляем саму карточку
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
func (r *Repository) GetCardByID(cardID uint) (models.Card, error) {
	var card models.Card
	if err := r.db.First(&card, cardID).Error; err != nil {
		return card, err
	}
	return card, nil
}
