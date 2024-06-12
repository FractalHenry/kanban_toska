package repository

import (
	"backend/models"
)

// Создание InformationalBlock
func (r *Repository) CreateInformationalBlock(block *models.InformationalBlock, userLogin string) error {
	if err := r.checkUserPermissionsForSpaceByBoardID(block.BoardID, userLogin, "createInformationalBlock"); err != nil {
		return err
	}

	return r.db.Create(&block).Error
}

// Удаление InformationalBlock по ID
func (r *Repository) DeleteInformationalBlockByBoardID(boardID uint, userLogin string) error {
	if err := r.checkUserPermissionsForSpaceByBoardID(boardID, userLogin, "deleteInformationalBlock"); err != nil {
		return err
	}

	return r.db.Where("board_id = ?", boardID).Delete(&models.InformationalBlock{}).Error
}

// Получение InformationalBlock по BoardID
func (r *Repository) GetInformationalBlockByBoardID(boardID uint) (models.InformationalBlock, error) {
	var block models.InformationalBlock
	err := r.db.Where("board_id = ?", boardID).First(&block).Error
	return block, err
}
