package repository

import (
	"backend/models"
	"fmt"
	"reflect"
)

func (r *Repository) CreateRoleOnSpace(roleOnSpace *models.RoleOnSpace, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnSpace.SpaceID, user.Login).First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner {
		return r.db.Create(roleOnSpace).Error
	} else if currentRole.IsAdmin {
		if roleOnSpace.IsOwner || roleOnSpace.IsAdmin {
			return fmt.Errorf("администратор не может создавать роли с правами владельца или администратора")
		}
		return r.db.Create(roleOnSpace).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для создания роли")
	}
}

func (r *Repository) CreateRoleOnBoard(roleOnBoard *models.RoleOnBoard, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("space_id = (SELECT space_id FROM boards WHERE board_id = ?) AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnBoard.SpaceID, user.Login).First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		return r.db.Create(roleOnBoard).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для создания роли для доски")
	}
}

func (r *Repository) UpdateRoleOnSpace(roleOnSpace *models.RoleOnSpace, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnSpace.SpaceID, user.Login).First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner {
		return r.db.Save(roleOnSpace).Error
	} else if currentRole.IsAdmin {
		if roleOnSpace.IsOwner {
			return fmt.Errorf("администратор не может изменять роль владельца")
		}
		return r.db.Save(roleOnSpace).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для изменения роли")
	}
}

func (r *Repository) UpdateRoleOnBoard(roleOnBoard *models.RoleOnBoard, boardIDs []uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}
	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = (SELECT space_id FROM role_on_boards WHERE role_on_board_id = ?) AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnBoard.RoleOnBoardID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		err = r.db.Save(roleOnBoard).Error
		if err != nil {
			return err
		}

		var updatedRoleOnBoard models.RoleOnBoard
		err = r.db.First(&updatedRoleOnBoard, roleOnBoard.RoleOnBoardID).Error
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
			if board.SpaceID == updatedRoleOnBoard.SpaceID {
				validBoardIDs = append(validBoardIDs, boardID)
			}
		}

		err = r.db.Where("role_on_board_id = ?", roleOnBoard.RoleOnBoardID).Delete(&models.BoardRoleOnBoard{}).Error
		if err != nil {
			return err
		}

		for _, boardID := range validBoardIDs {
			err = r.db.Create(&models.BoardRoleOnBoard{
				RoleOnBoardID: roleOnBoard.RoleOnBoardID,
				BoardID:       boardID,
			}).Error
			if err != nil {
				return err
			}
		}

		return nil
	} else {
		return fmt.Errorf("у пользователя нет прав для изменения роли для доски")
	}
}

func (r *Repository) AssociateUserWithRoleOnSpace(userLogin, targetUserLogin string, roleOnSpaceID uint) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}
	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("space_id = (SELECT space_id FROM role_on_spaces WHERE role_on_space_id = ?) AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnSpaceID, user.Login).First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner {
		return r.associateUserWithRole(targetUserLogin, roleOnSpaceID)
	} else if currentRole.IsAdmin {
		var roleToAdd models.RoleOnSpace
		err = r.db.First(&roleToAdd, roleOnSpaceID).Error
		if err != nil {
			return err
		}
		if roleToAdd.IsOwner {
			return fmt.Errorf("администратор не может добавлять роль владельца другому пользователю")
		}
		return r.associateUserWithRole(targetUserLogin, roleOnSpaceID)
	} else {
		return fmt.Errorf("у пользователя нет прав для добавления роли другому пользователю")
	}
}

func (r *Repository) AssociateUserWithRoleOnBoard(userLogin, targetUserLogin string, roleOnBoardID uint) error {
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
		return r.associateUserWithRole(targetUserLogin, roleOnBoardID)
	} else {
		return fmt.Errorf("у пользователя нет прав для добавления роли на доске другому пользователю")
	}
}

func (r *Repository) associateUserWithRole(userLogin string, roleID uint) error {
	targetUser, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	if roleID > 0 {
		roleType := reflect.TypeOf(roleID)
		if roleType == reflect.TypeOf(uint(0)) {
			return r.db.Model(&models.RoleOnSpace{}).Association("Users").Append(targetUser).Error
		} else if roleType == reflect.TypeOf(uint64(0)) {
			return r.db.Model(&models.RoleOnBoard{}).Association("Users").Append(targetUser).Error
		} else {
			return fmt.Errorf("неверный тип идентификатора роли")
		}
	}
	return fmt.Errorf("неверный идентификатор роли")
}

func (r *Repository) GetUserRolesInSpace(userLogin string, spaceID uint) ([]models.RoleOnSpace, error) {
	var roles []models.RoleOnSpace
	err := r.db.Model(&models.RoleOnSpace{}).
		Joins("JOIN user_role_on_spaces ON user_role_on_spaces.role_on_space_id = role_on_spaces.role_on_space_id").
		Where("user_role_on_spaces.login = ? AND role_on_spaces.space_id = ?", userLogin, spaceID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Repository) GetUserRolesOnBoard(userLogin string, boardID uint) ([]models.RoleOnBoard, error) {
	var roles []models.RoleOnBoard
	err := r.db.Model(&models.RoleOnBoard{}).
		Joins("JOIN board_role_on_boards ON board_role_on_boards.role_on_board_id = role_on_boards.role_on_board_id").
		Joins("JOIN boards ON boards.board_id = board_role_on_boards.board_id").
		Joins("JOIN user_board_role_on_boards ON user_board_role_on_boards.role_on_board_id = role_on_boards.role_on_board_id").
		Where("user_board_role_on_boards.login = ? AND boards.board_id = ?", userLogin, boardID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
