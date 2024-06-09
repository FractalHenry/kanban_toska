package repository

import (
	"backend/models"
	"fmt"
)

func (r *Repository) CreateSpace(space *models.Space) error {
	return r.db.Create(space).Error
}

func (r *Repository) CreateSpaceWithOwnerRole(spaceName string, userLogin string, roleOnSpaceName string) (*models.Space, *models.RoleOnSpace, error) {
	space := &models.Space{
		SpaceName: spaceName,
	}
	err := r.db.Create(space).Error
	if err != nil {
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
	err = r.db.Create(role).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.Model(&space).Association("RoleOnSpaces").Append(role).Error
	if err != nil {
		return nil, nil, err
	}

	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return nil, nil, err
	}
	err = r.db.Model(&role).Association("Users").Append(user).Error
	if err != nil {
		return nil, nil, err
	}

	return space, role, nil
}

func (r *Repository) UpdateSpaceName(spaceID uint, newSpaceName, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var space models.Space
	err = r.db.First(&space, spaceID).Error
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", space.SpaceID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner {
		space.SpaceName = newSpaceName
		return r.db.Save(&space).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для изменения имени пространства")
	}
}

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

func (r *Repository) GetSpace(spaceID uint) (*models.Space, error) {
	var space models.Space
	err := r.db.First(&space, spaceID).Error
	if err != nil {
		return nil, err
	}
	return &space, nil
}

func (r *Repository) GetSpaceByBoard(boardID uint) (*models.Space, error) {
	var board models.Board
	err := r.db.First(&board, boardID).Error
	if err != nil {
		return nil, err
	}

	var space models.Space
	err = r.db.First(&space, board.SpaceID).Error
	if err != nil {
		return nil, err
	}

	return &space, nil
}

func (r *Repository) ArchiveSpace(spaceID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var space models.Space
	err = r.db.First(&space, spaceID).Error
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", space.SpaceID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner {
		space.SpaceInArchive = true
		return r.db.Save(&space).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для перевода пространства в архив")
	}
}

func (r *Repository) UnarchiveSpace(spaceID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var space models.Space
	err = r.db.First(&space, spaceID).Error
	if err != nil {
		return err
	}

	var currentRole models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", space.SpaceID, user.Login).
		First(&currentRole).Error
	if err != nil {
		return err
	}

	if currentRole.IsOwner {
		space.SpaceInArchive = false
		return r.db.Save(&space).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для восстановления пространства из архива")
	}
}
