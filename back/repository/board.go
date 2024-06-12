package repository

import (
	"backend/models"
	"fmt"
)

// Создает новую доску
func (r *Repository) CreateBoard(board *models.Board, userLogin string) error {
	if err := r.checkUserPermissionsForSpace(board.SpaceID, userLogin, "createBoard"); err != nil {
		return err
	}
	return r.db.Create(board).Error
}

// Связывает роль с досками
func (r *Repository) AssociateRoleWithBoards(roleOnBoardID uint, boardIDs []uint, userLogin string) error {
	if err := r.checkUserPermissionsForRole(roleOnBoardID, userLogin, "associateRoleWithBoards"); err != nil {
		return err
	}

	var roleOnBoard models.RoleOnBoard
	if err := r.db.First(&roleOnBoard, roleOnBoardID).Error; err != nil {
		return err
	}

	validBoardIDs := make([]uint, 0)
	for _, boardID := range boardIDs {
		var board models.Board
		if err := r.db.First(&board, boardID).Error; err != nil {
			return err
		}
		if board.SpaceID == roleOnBoard.SpaceID {
			validBoardIDs = append(validBoardIDs, boardID)
		}
	}

	for _, boardID := range validBoardIDs {
		err := r.db.Create(&models.BoardRoleOnBoard{
			RoleOnBoardID: roleOnBoardID,
			BoardID:       boardID,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// Обновляет имя доски
func (r *Repository) UpdateBoardName(boardID uint, newBoardName, userLogin string) error {
	if err := r.checkUserPermissionsForSpaceByBoardID(boardID, userLogin, "updateBoardName"); err != nil {
		return err
	}

	var board models.Board
	if err := r.db.First(&board, boardID).Error; err != nil {
		return err
	}

	board.BoardName = newBoardName
	return r.db.Save(&board).Error
}

// Возвращает доски пользователя
func (r *Repository) GetUserBoards(userLogin string) ([]models.Board, error) {
	var boards []models.Board

	err := r.db.Table("boards").
		Joins("JOIN board_role_on_boards br ON br.board_id = boards.board_id").
		Joins("JOIN user_board_role_on_boards ubr ON ubr.board_role_on_board_id = br.role_on_board_id").
		Where("ubr.login = ?", userLogin).
		Find(&boards).Error

	if err != nil {
		return nil, err
	}
	return boards, nil
}

// Возвращает доску по ее ID
func (r *Repository) GetBoard(boardID uint) (*models.Board, error) {
	var board models.Board
	err := r.db.First(&board, boardID).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

// Переводит доску в архив
func (r *Repository) ArchiveBoard(boardID uint, userLogin string) error {
	if err := r.checkUserPermissionsForSpaceByBoardID(boardID, userLogin, "archiveBoard"); err != nil {
		return err
	}

	var board models.Board
	if err := r.db.First(&board, boardID).Error; err != nil {
		return err
	}

	board.BoardInArchive = true
	return r.db.Save(&board).Error
}

// Восстанавливает доску из архива
func (r *Repository) UnarchiveBoard(boardID uint, userLogin string) error {
	if err := r.checkUserPermissionsForSpaceByBoardID(boardID, userLogin, "unarchiveBoard"); err != nil {
		return err
	}

	var board models.Board
	if err := r.db.First(&board, boardID).Error; err != nil {
		return err
	}

	board.BoardInArchive = false
	return r.db.Save(&board).Error
}

// Проверяет права пользователя на выполнение действия в пространстве
func (r *Repository) checkUserPermissionsForSpace(spaceID uint, userLogin, action string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", spaceID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	switch action {
	case "createBoard", "associateRoleWithBoards":
		if currentRole.IsOwner || currentRole.IsAdmin || currentRole.CanEdit {
			return nil
		}
	case "updateBoardName", "archiveBoard", "unarchiveBoard":
		if currentRole.IsOwner || currentRole.IsAdmin {
			return nil
		}
	}

	return fmt.Errorf("у пользователя нет прав для выполнения действия: %s", action)
}

// Проверяет права пользователя на выполнение действия для роли
func (r *Repository) checkUserPermissionsForRole(roleOnBoardID uint, userLogin, action string) error {
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
		return nil
	}

	return fmt.Errorf("у пользователя нет прав для выполнения действия: %s", action)
}

// Проверяет права пользователя на выполнение действия в пространстве по ID доски
func (r *Repository) checkUserPermissionsForSpaceByBoardID(boardID uint, userLogin, action string) error {
	var board models.Board
	if err := r.db.First(&board, boardID).Error; err != nil {
		return err
	}

	return r.checkUserPermissionsForSpace(board.SpaceID, userLogin, action)
}

func (r *Repository) GetBoardDetails(boardID uint) (*models.Board, []models.Card, []models.Task, *models.InformationalBlock, error) {
	var board models.Board
	if err := r.db.First(&board, boardID).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	var cards []models.Card
	if err := r.db.Where("board_id = ?", boardID).Find(&cards).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	var tasks []models.Task
	for _, card := range cards {
		cardTasks, err := r.GetCardTasks(card.CardID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		tasks = append(tasks, cardTasks...)
	}

	infoBlock, err := r.GetInformationalBlockByBoardID(boardID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return &board, cards, tasks, infoBlock, nil
}

func (r *Repository) GetBoardUsers(boardID uint) ([]struct {
	Login       string
	RoleOnBoard models.RoleOnBoard
}, error) {
	var users []struct {
		Login       string
		RoleOnBoard models.RoleOnBoard
	}

	err := r.db.Table("user_board_role_on_boards").
		Select("users.login, role_on_boards.role_on_board_id, role_on_boards.role_on_board_name, role_on_boards.space_id").
		Joins("JOIN users ON users.login = user_board_role_on_boards.login").
		Joins("JOIN role_on_boards ON role_on_boards.role_on_board_id = user_board_role_on_boards.role_on_board_id").
		Where("user_board_role_on_boards.board_id = ?", boardID).
		Scan(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
