package repository

import (
	"backend/models"
	"fmt"

	"gorm.io/gorm"
)

// Функция создания карточки
func (r *Repository) CreateCard(card *models.Card, userLogin string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Получаем доску, к которой относится карточка
		var board models.Board
		if err := tx.Where("board_id = ?", card.BoardID).First(&board).Error; err != nil {
			return err
		}

		// Получаем пространство, к которому относится доска
		var space models.Space
		if err := tx.Where("space_id = ?", board.SpaceID).First(&space).Error; err != nil {
			return err
		}

		// Проверяем роль пользователя на доске
		var userBoardRole models.UserBoardRoleOnBoard
		if err := tx.Where("login = ? AND board_id = ?", userLogin, card.BoardID).First(&userBoardRole).Error; err == nil {
			// Роль на доске найдена
			var boardRole models.BoardRoleOnBoard
			if err := tx.Where("role_on_board_id = ?", userBoardRole.RoleOnBoardID).First(&boardRole).Error; err != nil {
				return err
			}

			// Проверяем, есть ли у пользователя права на редактирование доски
			if boardRole.CanEdit {
				return tx.Create(card).Error
			}

			return fmt.Errorf("у пользователя нет прав для создания карточки на этой доске")
		} else if err != gorm.ErrRecordNotFound {
			// Произошла ошибка, не связанная с отсутствием записи
			return err
		}

		// Если роль на доске не найдена, проверяем роль пользователя в пространстве
		var userRoleOnSpace models.UserRoleOnSpace
		if err := tx.Where("login = ? AND space_id = ?", userLogin, board.SpaceID).First(&userRoleOnSpace).Error; err == nil {
			// Роль на пространстве найдена
			var roleOnSpace models.RoleOnSpace
			if err := tx.Where("role_on_space_id = ?", userRoleOnSpace.RoleOnSpaceID).First(&roleOnSpace).Error; err != nil {
				return err
			}

			// Проверяем, есть ли у пользователя права на редактирование пространства
			if roleOnSpace.CanEdit || roleOnSpace.IsAdmin || roleOnSpace.IsOwner {
				return tx.Create(card).Error
			}

			return fmt.Errorf("у пользователя нет прав для создания карточки в этом пространстве")
		} else if err != gorm.ErrRecordNotFound {
			// Произошла ошибка, не связанная с отсутствием записи
			return err
		}

		// Если роль пользователя не найдена ни на доске, ни в пространстве
		return fmt.Errorf("у пользователя нет прав для создания карточки")
	})
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
