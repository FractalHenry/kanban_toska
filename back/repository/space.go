package repository

import (
	"backend/models"
	"fmt"
)

// Cоздает новое пространство
func (r *Repository) CreateSpace(space *models.Space) error {
	return r.db.Create(space).Error
}

// Создает пространство с ролью владельца
func (r *Repository) CreateSpaceWithOwnerRole(spaceName, userLogin, roleOnSpaceName string) (*models.Space, *models.RoleOnSpace, error) {
	space := &models.Space{
		SpaceName: spaceName,
	}
	if err := r.db.Create(space).Error; err != nil {
		return nil, nil, err
	}

	if roleOnSpaceName == "" {
		roleOnSpaceName = "Владелец"
	}

	role := &models.RoleOnSpace{
		RoleOnSpaceName: roleOnSpaceName,
		IsOwner:         true,
		SpaceID:         space.SpaceID,
	}

	if err := r.db.Create(&role).Error; err != nil {
		return nil, nil, err
	}

	if err := r.associateUserWithRole(userLogin, role.RoleOnSpaceID, "space"); err != nil {
		return nil, nil, err
	}

	return space, role, nil
}

// Обновляет имя пространства
func (r *Repository) UpdateSpaceName(spaceID uint, newSpaceName, userLogin string) error {
	if err := r.checkUserOwnership(spaceID, userLogin); err != nil {
		return err
	}

	var space models.Space
	if err := r.db.First(&space, spaceID).Error; err != nil {
		return err
	}

	space.SpaceName = newSpaceName
	return r.db.Save(&space).Error
}

// Возвращает пространства пользователя
func (r *Repository) GetUserSpaces(userLogin string) ([]models.Space, error) {
	var spaces []models.Space
	err := r.db.Model(&models.Space{}).
		Joins("JOIN role_on_spaces ON role_on_spaces.space_id = spaces.space_id").
		Joins("JOIN user_role_on_spaces ON user_role_on_spaces.role_on_space_id = role_on_spaces.role_on_space_id").
		Where("user_role_on_spaces.login = ?", userLogin).
		Find(&spaces).Error
	if err != nil {
		return nil, err
	}
	return spaces, nil
}

// Возвращает пространство по его ID
func (r *Repository) GetSpace(spaceID uint) (*models.Space, error) {
	var space models.Space
	err := r.db.First(&space, spaceID).Error
	if err != nil {
		return nil, err
	}
	return &space, nil
}

// Возвращает пространство по ID доски
func (r *Repository) GetSpaceByBoard(boardID uint) (*models.Space, error) {
	var board models.Board
	if err := r.db.First(&board, boardID).Error; err != nil {
		return nil, err
	}

	var space models.Space
	if err := r.db.First(&space, board.SpaceID).Error; err != nil {
		return nil, err
	}

	return &space, nil
}

// Переводит пространство в архив
func (r *Repository) ArchiveSpace(spaceID uint, userLogin string) error {
	if err := r.checkUserOwnership(spaceID, userLogin); err != nil {
		return err
	}

	var space models.Space
	if err := r.db.First(&space, spaceID).Error; err != nil {
		return err
	}

	space.SpaceInArchive = true
	return r.db.Save(&space).Error
}

// Восстанавливает пространство из архива
func (r *Repository) UnarchiveSpace(spaceID uint, userLogin string) error {
	if err := r.checkUserOwnership(spaceID, userLogin); err != nil {
		return err
	}

	var space models.Space
	if err := r.db.First(&space, spaceID).Error; err != nil {
		return err
	}

	space.SpaceInArchive = false
	return r.db.Save(&space).Error
}

// Проверяет, является ли пользователь владельцем пространства
func (r *Repository) checkUserOwnership(spaceID uint, userLogin string) error {
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

	if !currentRole.IsOwner {
		return fmt.Errorf("у пользователя нет прав для выполнения этого действия")
	}

	return nil
}

func (r *Repository) GetSpaceOwner(spaceID uint) (string, error) {
	var roleOnSpace models.RoleOnSpace
	err := r.db.Where("space_id = ? AND is_owner = true", spaceID).First(&roleOnSpace).Error
	if err != nil {
		return "", err
	}

	var userRoleOnSpace models.UserRoleOnSpace
	err = r.db.Where("role_on_space_id = ?", roleOnSpace.RoleOnSpaceID).First(&userRoleOnSpace).Error
	if err != nil {
		return "", err
	}

	var user models.User
	err = r.db.Where("login = ?", userRoleOnSpace.Login).First(&user).Error
	if err != nil {
		return "", err
	}

	return user.Login, nil
}

func (r *Repository) GetSpaceBoards(spaceID uint) ([]models.Board, error) {
	var boards []models.Board
	err := r.db.Where("space_id = ?", spaceID).Find(&boards).Error
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (r *Repository) GetSpaceUsers(spaceID uint) (*[]models.UserRoleOnSpace, error) {
	var userRoleOnSpaces *[]models.UserRoleOnSpace
	err := r.db.Where("role_on_spaces.space_id = ?", spaceID).
		Joins("JOIN role_on_spaces ON user_role_on_spaces.role_on_space_id = role_on_spaces.role_on_space_id").
		Joins("JOIN users ON user_role_on_spaces.login = users.login").
		Select("users.login").
		Find(&userRoleOnSpaces).Error
	if err != nil {
		return nil, err
	}

	return userRoleOnSpaces, nil
}
