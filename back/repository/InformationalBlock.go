package repository

import (
	"backend/models"
	"fmt"
)

// Создание InformationalBlock
func (r *Repository) CreateInformationalBlock(block *models.InformationalBlock, userLogin string) error {

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRolesForBlock(block, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Create(&block).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для создания цвета задания")
	}

}

// Удаление InformationalBlock по ID
func (r *Repository) DeleteInformationalBlockByBoardID(boardID uint, userLogin string) error {
	block, err := r.GetInformationalBlockByBoardID(boardID)

	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRolesForBlock(block, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(roleOnSpace, boardRole) {
		return r.db.Where("board_id = ?", boardID).Delete(&models.InformationalBlock{}).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления информационного блока")
	}
}

// Получение InformationalBlock по BoardID
func (r *Repository) GetInformationalBlockByBoardID(boardID uint) (*models.InformationalBlock, error) {

	var block *models.InformationalBlock
	err := r.db.Where("board_id = ?", boardID).First(&block).Error

	if err != nil {
		return nil, err
	}

	return block, nil
}

func (r *Repository) findBoardAndRolesForBlock(input *models.InformationalBlock, userLogin string) (models.User, models.Board, models.RoleOnSpace, models.BoardRoleOnBoard, error) {
	userPtr, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}
	user := *userPtr

	var board models.Board
	if err := r.db.First(&board, input.BoardID).Error; err != nil {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}

	var roleOnSpace models.RoleOnSpace
	if err := r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&roleOnSpace).Error; err != nil {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}

	var boardRole models.BoardRoleOnBoard
	if err := r.db.Model(&models.BoardRoleOnBoard{}).
		Where("board_id = ? AND role_on_board_id IN (SELECT role_on_board_id FROM user_board_role_on_boards WHERE login = ?)", board.BoardID, user.Login).
		First(&boardRole).Error; err != nil {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}

	return user, board, roleOnSpace, boardRole, nil
}
