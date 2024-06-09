package repository

import (
	"backend/models"
	"fmt"
)

func (r *Repository) CreateBoard(board *models.Board, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner || currentRole.IsAdmin || currentRole.CanEdit {
		return r.db.Create(board).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для создания доски")
	}
}

func (r *Repository) AssociateRoleWithBoards(roleOnBoardID uint, boardIDs []uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = (SELECT space_id FROM role_on_boards WHERE role_on_board_id = ?) AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnBoardID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		var roleOnBoard models.RoleOnBoard
		err = r.db.First(&roleOnBoard, roleOnBoardID).Error
		if err != nil {
			return err
		}

		validBoardIDs := make([]uint, 0)
		for _, boardID := range boardIDs {
			var board models.Board
			err = r.db.First(&board, boardID).Error
			if err != nil {
				return err
			}
			if board.SpaceID == roleOnBoard.SpaceID {
				validBoardIDs = append(validBoardIDs, boardID)
			}
		}

		for _, boardID := range validBoardIDs {
			err = r.db.Create(&models.BoardRoleOnBoard{
				RoleOnBoardID: roleOnBoardID,
				BoardID:       boardID,
			}).Error
			if err != nil {
				return err
			}
		}

		return nil
	} else {
		return fmt.Errorf("у пользователя нет прав для связывания роли с досками")
	}
}

func (r *Repository) UpdateBoardName(boardID uint, newBoardName, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var board models.Board
	err = r.db.First(&board, boardID).Error
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		board.BoardName = newBoardName
		return r.db.Save(&board).Error
	} else {
		// Проверка права на редактирование доски
		// var roleOnBoard models.RoleOnBoard
		// err = r.db.Model(&models.RoleOnBoard{}).
		//     Where("role_on_board_id IN (SELECT role_on_board_id FROM user_board_role_on_boards WHERE login = ?) AND can_edit = true", user.Login).
		//     First(&roleOnBoard).Error
		// if err != nil {
		//     return fmt.Errorf("у пользователя нет прав для изменения имени доски")
		// }
		//board.BoardName = newBoardName
		//return r.db.Save(&board).Error
		return fmt.Errorf("у пользователя нет прав для изменения имени доски")
	}
}

func (r *Repository) GetUserBoards(userLogin string) ([]models.Board, error) {
	var boards []models.Board
	err := r.db.Model(&models.Board{}).
		Joins("JOIN board_role_on_boards ON board_role_on_boards.board_id = boards.board_id").
		Joins("JOIN role_on_boards ON role_on_boards.role_on_board_id = board_role_on_boards.role_on_board_id").
		Joins("JOIN user_board_role_on_boards ON user_board_role_on_boards.role_on_board_id = role_on_boards.role_on_board_id").
		Where("user_board_role_on_boards.login = ?", userLogin).
		Find(&boards).Error
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (r *Repository) GetBoard(boardID uint) (*models.Board, error) {
	var board models.Board
	err := r.db.First(&board, boardID).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *Repository) ArchiveBoard(boardID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var board models.Board
	err = r.db.First(&board, boardID).Error
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		board.BoardInArchive = true
		return r.db.Save(&board).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для перевода доски в архив")
	}
}

func (r *Repository) UnarchiveBoard(boardID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var board models.Board
	err = r.db.First(&board, boardID).Error
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		board.BoardInArchive = false
		return r.db.Save(&board).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для восстановления доски из архива")
	}
}
