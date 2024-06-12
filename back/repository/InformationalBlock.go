package repository

import (
	"backend/models"
	"fmt"

	"gorm.io/gorm"
)

// Создание InformationalBlock
func (r *Repository) CreateInformationalBlock(block *models.InformationalBlock, userLogin string) error {
	_, _, roleOnSpace, boardRole, err := r.findBoardAndRolesForBlock(block.BoardID, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(*roleOnSpace, *boardRole) {
		return r.db.Create(&block).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для создания информационного блока")
	}
}

// Удаление InformationalBlock по ID
func (r *Repository) DeleteInformationalBlockByBoardID(blockID uint, userLogin string) error {
	block, err := r.GetInformationalBlockByID(blockID)
	if err != nil {
		return err
	}

	_, _, roleOnSpace, boardRole, err := r.findBoardAndRolesForBlock(block.BoardID, userLogin)
	if err != nil {
		return err
	}

	if r.hasEditPermissions(*roleOnSpace, *boardRole) {
		return r.db.Where("informational_block_id = ?", blockID).Delete(&models.InformationalBlock{}).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления информационного блока")
	}
}

// Получение InformationalBlock по BoardID
func (r *Repository) GetInformationalBlockByID(blockID uint) (*models.InformationalBlock, error) {
	var block models.InformationalBlock
	err := r.db.Where("informational_block_id = ?", blockID).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *Repository) GetInformationalBlocksByBoardID(boardID uint) ([]models.InformationalBlock, error) {
	var blocks []models.InformationalBlock
	err := r.db.Where("board_id = ?", boardID).Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (r *Repository) findBoardAndRolesForBlock(boardID uint, userLogin string) (*models.User, *models.Board, *models.RoleOnSpace, *models.BoardRoleOnBoard, error) {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var board *models.Board
	if err := r.db.First(&board, boardID).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	var roleOnSpace *models.RoleOnSpace
	if err := r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&roleOnSpace).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, nil, nil, err
	}

	var boardRole *models.BoardRoleOnBoard
	if err := r.db.Model(&models.BoardRoleOnBoard{}).
		Where("board_id = ? AND role_on_board_id IN (SELECT role_on_board_id FROM user_board_role_on_boards WHERE login = ?)", boardID, user.Login).
		First(&boardRole).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, nil, nil, err
	}

	return user, board, roleOnSpace, boardRole, nil
}
