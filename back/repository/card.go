package repository

import (
	"backend/models"
	"fmt"
)

func (r *Repository) CreateCard(card *models.Card, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var board models.Board
	err = r.db.First(&board, card.BoardID).Error
	if err != nil {
		return err
	}

	var currentRoleOnSpace models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&currentRoleOnSpace).Error
	if err != nil {
		return err
	}

	var currentBoardRole models.BoardRoleOnBoard
	err = r.db.Model(&models.BoardRoleOnBoard{}).
		Where("board_id = ? AND role_on_board_id IN (SELECT role_on_board_id FROM user_board_role_on_boards WHERE login = ?)", board.BoardID, user.Login).
		First(&currentBoardRole).Error
	if err != nil {
		return err
	}

	if currentRoleOnSpace.IsOwner || currentRoleOnSpace.IsAdmin || currentRoleOnSpace.CanEdit || currentBoardRole.CanEdit {
		return r.db.Create(card).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для создания карточки на этой доске")
	}
}
