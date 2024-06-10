package repository

import (
	"backend/models"
	"fmt"
)

func (r *Repository) CreateRoleOnSpace(roleOnSpace *models.RoleOnSpace, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnSpace.SpaceID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if roleOnSpace.IsOwner {
		return fmt.Errorf("роль владельца не может быть создана")
	}

	if currentRole.IsOwner {
		return r.db.Create(roleOnSpace).Error
	} else if currentRole.IsAdmin {
		if roleOnSpace.IsAdmin {
			return fmt.Errorf("администратор не может создавать роли с правами администратора")
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

	var existingRole models.RoleOnSpace
	err = r.db.First(&existingRole, roleOnSpace.RoleOnSpaceID).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner {
		if existingRole.IsOwner && !roleOnSpace.IsOwner {
			return fmt.Errorf("нельзя убрать права владельца с роли владельца")
		}
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
		return fmt.Errorf("не удалось найти пользователя: %v", err)
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", user.Login).First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль пользователя: %v", err)
	}

	var roleToAdd models.RoleOnSpace
	err = r.db.First(&roleToAdd, roleOnSpaceID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль: %v", err)
	}
	if roleToAdd.IsOwner {
		return fmt.Errorf("роль владельца не может быть привязана к пользователю")
	}

	if currentRole.IsOwner {
		return r.associateUserWithRole(targetUserLogin, roleOnSpaceID, "space")
	} else if currentRole.IsAdmin {
		if roleToAdd.IsAdmin {
			return fmt.Errorf("администратор не может привязывать роль администратора")
		}
		return r.associateUserWithRole(targetUserLogin, roleOnSpaceID, "space")
	} else {
		return fmt.Errorf("у пользователя нет прав для привязки роли другому пользователю")
	}
}

func (r *Repository) AssociateUserWithRoleOnBoard(userLogin, targetUserLogin string, roleOnBoardID uint) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя: %v", err)
	}

	// Проверка текущей роли пользователя
	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = (SELECT space_id FROM role_on_boards WHERE role_on_board_id = ?) AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnBoardID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль пользователя: %v", err)
	}

	// Проверка роли, которую пытаются привязать
	var roleOnBoard models.RoleOnBoard
	err = r.db.First(&roleOnBoard, roleOnBoardID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль: %v", err)
	}
	var roleInSpace models.RoleOnSpace
	err = r.db.First(&roleInSpace, roleOnBoard.SpaceID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль в пространстве: %v", err)
	}
	if roleInSpace.IsOwner {
		return fmt.Errorf("роль владельца не может быть привязана к пользователю")
	}

	// Проверка прав текущего пользователя
	if currentRole.IsOwner {
		return r.associateUserWithRole(targetUserLogin, roleOnBoardID, "board")
	} else if currentRole.IsAdmin {
		if roleInSpace.IsAdmin {
			return fmt.Errorf("администратор не может привязывать роль администратора")
		}
		return r.associateUserWithRole(targetUserLogin, roleOnBoardID, "board")
	} else {
		return fmt.Errorf("у пользователя нет прав для привязки роли другому пользователю")
	}
}

func (r *Repository) associateUserWithRole(userLogin string, roleID uint, roleType string) error {
	targetUser, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	if roleID > 0 {
		switch roleType {
		case "space":
			userRoleOnSpace := models.UserRoleOnSpace{
				RoleOnSpaceID: roleID,
				Login:         targetUser.Login,
			}
			return r.db.Create(&userRoleOnSpace).Error
		case "board":
			userBoardRoleOnBoard := models.UserBoardRoleOnBoard{
				BoardRoleOnBoardID: roleID,
				Login:              targetUser.Login,
			}
			return r.db.Create(&userBoardRoleOnBoard).Error
		default:
			return fmt.Errorf("неверный тип роли")
		}
	}
	return fmt.Errorf("неверный идентификатор роли")
}

func (r *Repository) AssociateRoleWithBoard(userLogin string, boardRoleOnBoard *models.BoardRoleOnBoard, boardID uint) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя: %v", err)
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Joins("JOIN boards ON boards.space_id = role_on_spaces.space_id").
		Where("boards.board_id = ? AND role_on_spaces.role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", boardID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль пользователя в пространстве: %v", err)
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		err = r.db.Create(boardRoleOnBoard).Error
		if err != nil {
			return fmt.Errorf("не удалось ассоциировать роль с доской: %v", err)
		}

		return nil
	}

	return fmt.Errorf("у пользователя нет прав для связывания роли с доской")
}

func (r *Repository) DeleteRoleOnSpace(roleOnSpaceID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя: %v", err)
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", user.Login).First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль пользователя: %v", err)
	}

	var roleToDelete models.RoleOnSpace
	err = r.db.First(&roleToDelete, roleOnSpaceID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль: %v", err)
	}

	if roleToDelete.IsOwner {
		return fmt.Errorf("роль владельца не может быть удалена")
	}

	if currentRole.IsOwner {
		return r.db.Delete(&models.RoleOnSpace{}, roleOnSpaceID).Error
	} else if currentRole.IsAdmin {
		if roleToDelete.IsAdmin {
			return fmt.Errorf("администратор не может удалять роль администратора")
		}
		return r.db.Delete(&models.RoleOnSpace{}, roleOnSpaceID).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления роли")
	}
}

func (r *Repository) DeleteRoleOnBoard(roleOnBoardID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя: %v", err)
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", user.Login).First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль пользователя: %v", err)
	}

	var roleToDelete models.RoleOnBoard
	err = r.db.First(&roleToDelete, roleOnBoardID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль: %v", err)
	}

	var roleInSpace models.RoleOnSpace
	err = r.db.First(&roleInSpace, roleToDelete.SpaceID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль в пространстве: %v", err)
	}
	if roleInSpace.IsOwner {
		return fmt.Errorf("роль владельца не может быть удалена")
	}

	if currentRole.IsOwner {
		return r.db.Delete(&models.RoleOnBoard{}, roleOnBoardID).Error
	} else if currentRole.IsAdmin {
		if roleInSpace.IsAdmin {
			return fmt.Errorf("администратор не может удалять роль администратора")
		}
		return r.db.Delete(&models.RoleOnBoard{}, roleOnBoardID).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления роли")
	}
}

func (r *Repository) DeleteBoardRoleOnBoardAssociation(boardRoleOnBoardID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя: %v", err)
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Joins("JOIN boards ON boards.space_id = role_on_spaces.space_id").
		Where("boards.board_id = (SELECT board_id FROM board_role_on_boards WHERE board_role_on_board_id = ?) AND role_on_spaces.role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", boardRoleOnBoardID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль пользователя в пространстве: %v", err)
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		return r.db.Delete(&models.BoardRoleOnBoard{}, boardRoleOnBoardID).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления связи между ролью и доской")
	}
}

func (r *Repository) DeleteUserRoleOnSpaceAssociation(userRoleOnSpaceID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя: %v", err)
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).Where("role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", user.Login).First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль пользователя: %v", err)
	}

	var roleToDelete models.UserRoleOnSpace
	err = r.db.First(&roleToDelete, userRoleOnSpaceID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти связь роли: %v", err)
	}

	var roleInSpace models.RoleOnSpace
	err = r.db.First(&roleInSpace, roleToDelete.RoleOnSpaceID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль в пространстве: %v", err)
	}
	if roleInSpace.IsOwner {
		return fmt.Errorf("роль владельца не может быть удалена")
	}

	if currentRole.IsOwner {
		return r.db.Delete(&models.UserRoleOnSpace{}, userRoleOnSpaceID).Error
	} else if currentRole.IsAdmin {
		if roleInSpace.IsAdmin {
			return fmt.Errorf("администратор не может удалять связь роли администратора")
		}
		return r.db.Delete(&models.UserRoleOnSpace{}, userRoleOnSpaceID).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления связи роли")
	}
}

func (r *Repository) DeleteUserBoardRoleOnBoardAssociation(userBoardRoleOnBoardID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя: %v", err)
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Joins("JOIN boards ON boards.space_id = role_on_spaces.space_id").
		Where("boards.board_id = (SELECT board_id FROM user_board_role_on_boards WHERE user_board_role_on_board_id = ?) AND role_on_spaces.role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", userBoardRoleOnBoardID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль пользователя в пространстве: %v", err)
	}

	if currentRole.IsOwner || currentRole.IsAdmin {
		return r.db.Delete(&models.UserBoardRoleOnBoard{}, userBoardRoleOnBoardID).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления связи между пользователем и ролью доски")
	}
}

// Метод передачи роли владельца от одного пользователя другому
func (r *Repository) TransferOwnership(userLogin, newOwnerLogin string, roleOnSpaceID uint) error {
	// Поиск текущего пользователя
	currentUser, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти текущего пользователя: %v", err)
	}

	// Проверка текущей роли пользователя
	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("role_on_space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", roleOnSpaceID, currentUser.Login).
		First(&currentRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль текущего пользователя: %v", err)
	}

	// Проверка наличия прав у текущего пользователя для передачи роли владельца
	if !currentRole.IsOwner {
		return fmt.Errorf("только владелец может передавать права владельца")
	}

	// Поиск нового пользователя
	newOwner, err := r.FindUserByLogin(newOwnerLogin)
	if err != nil {
		return fmt.Errorf("не удалось найти нового пользователя: %v", err)
	}

	// Поиск текущей роли нового пользователя в данном пространстве
	var newOwnerRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", currentRole.SpaceID, newOwner.Login).
		First(&newOwnerRole).Error
	if err != nil {
		return fmt.Errorf("не удалось найти роль нового пользователя: %v", err)
	}

	// Проверка, не является ли новый пользователь уже владельцем
	if newOwnerRole.IsOwner {
		return fmt.Errorf("новый пользователь уже является владельцем")
	}

	// Начало транзакции
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("не удалось начать транзакцию: %v", tx.Error)
	}

	// Обновление роли текущего пользователя, убирая права владельца
	currentRole.IsOwner = false
	err = tx.Save(&currentRole).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("не удалось обновить роль текущего пользователя: %v", err)
	}

	// Назначение роли владельца новому пользователю
	newOwnerRole.IsOwner = true
	err = tx.Save(&newOwnerRole).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("не удалось назначить роль владельца новому пользователю: %v", err)
	}

	// Коммит транзакции
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("не удалось завершить транзакцию: %v", err)
	}

	return nil
}
