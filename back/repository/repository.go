package repository

import (
	"backend/models"
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Функция для нахождения пользователя, доски и ролей
func (r *Repository) findBoardAndRoles(input interface{}, userLogin string) (models.User, models.Board, models.RoleOnSpace, models.BoardRoleOnBoard, error) {
	userPtr, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}
	user := *userPtr

	var cardID uint
	switch v := input.(type) {
	case *models.Task:
		cardID = v.CardID
	case *models.Card:
		cardID = v.CardID
	case *models.Checklist:
		cardID = v.Task.CardID
	case *models.ChecklistElement:
		cardID = v.Checklist.Task.CardID
	default:
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, errors.New("unsupported input type")
	}

	var card models.Card
	if err := r.db.First(&card, cardID).Error; err != nil {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}

	var board models.Board
	if err := r.db.First(&board, card.BoardID).Error; err != nil {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}

	var roleOnSpace models.RoleOnSpace = models.RoleOnSpace{}
	if err := r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&roleOnSpace).Error; err != nil && err != gorm.ErrRecordNotFound {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}

	var boardRole models.BoardRoleOnBoard = models.BoardRoleOnBoard{}
	if err := r.db.Model(&models.BoardRoleOnBoard{}).
		Where("board_id = ? AND role_on_board_id IN (SELECT role_on_board_id FROM user_board_role_on_boards WHERE login = ?)", board.BoardID, user.Login).
		First(&boardRole).Error; err != nil && err != gorm.ErrRecordNotFound {
		return models.User{}, models.Board{}, models.RoleOnSpace{}, models.BoardRoleOnBoard{}, err
	}

	return user, board, roleOnSpace, boardRole, nil
}

// Функция проверки прав на редактирование
func (r *Repository) hasEditPermissions(roleOnSpace models.RoleOnSpace, boardRole models.BoardRoleOnBoard) bool {
	return roleOnSpace.IsOwner || roleOnSpace.IsAdmin || roleOnSpace.CanEdit || boardRole.CanEdit
}
