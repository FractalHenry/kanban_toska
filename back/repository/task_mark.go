package repository

import (
	"backend/models"
	"fmt"
)

// Проверяет права пользователя на выполнение действия для задачи
func (r *Repository) checkUserPermissionsForTask(taskID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var task models.Task
	if err := r.db.First(&task, taskID).Error; err != nil {
		return err
	}

	var card models.Card
	if err := r.db.First(&card, task.CardID).Error; err != nil {
		return err
	}

	var board models.Board
	if err := r.db.First(&board, card.BoardID).Error; err != nil {
		return err
	}

	var currentRoleOnSpace models.RoleOnSpace
	roleOnSpaceErr := r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&currentRoleOnSpace).Error

	var currentBoardRole models.BoardRoleOnBoard
	boardRoleErr := r.db.Model(&models.BoardRoleOnBoard{}).
		Where("board_id = ? AND role_on_board_id IN (SELECT role_on_board_id FROM user_board_role_on_boards WHERE login = ?)", board.BoardID, user.Login).
		First(&currentBoardRole).Error

	if roleOnSpaceErr != nil && boardRoleErr != nil {
		return fmt.Errorf("у пользователя нет прав для выполнения действия")
	}

	if currentRoleOnSpace.IsOwner || currentRoleOnSpace.IsAdmin || currentRoleOnSpace.CanEdit || currentBoardRole.CanEdit {
		return nil
	}

	return fmt.Errorf("у пользователя нет прав для выполнения действия")
}

// Создает новую метку
func (r *Repository) CreateMark(mark *models.Mark, userLogin string) error {
	if err := r.checkUserPermissionsForTask(mark.TaskID, userLogin); err != nil {
		return err
	}
	return r.db.Create(mark).Error
}

// Обновляет метку
func (r *Repository) UpdateMark(mark *models.Mark, userLogin string) error {
	if err := r.checkUserPermissionsForTask(mark.TaskID, userLogin); err != nil {
		return err
	}
	return r.db.Save(mark).Error
}

// Удаляет метку
func (r *Repository) DeleteMark(markID uint, userLogin string) error {
	var mark models.Mark
	if err := r.db.First(&mark, markID).Error; err != nil {
		return err
	}
	if err := r.checkUserPermissionsForTask(mark.TaskID, userLogin); err != nil {
		return err
	}
	return r.db.Delete(&mark).Error
}

// Создает новое имя метки
func (r *Repository) CreateMarkName(markName *models.MarkName, userLogin string) error {
	var task models.Task
	if err := r.db.First(&task, markName.Mark.TaskID).Error; err != nil {
		return err
	}
	if err := r.checkUserPermissionsForTask(task.TaskID, userLogin); err != nil {
		return err
	}
	return r.db.Create(markName).Error
}

// Обновляет имя метки
func (r *Repository) UpdateMarkName(markName *models.MarkName, userLogin string) error {
	var task models.Task
	if err := r.db.First(&task, markName.Mark.TaskID).Error; err != nil {
		return err
	}
	if err := r.checkUserPermissionsForTask(task.TaskID, userLogin); err != nil {
		return err
	}
	return r.db.Save(markName).Error
}

func (r *Repository) GetTaskMarkNamesAndMarks(taskID uint) (*[]models.MarkName, *[]models.Mark, error) {
	var markNames []models.MarkName
	var marks []models.Mark

	err := r.db.Model(&models.MarkName{}).
		Preload("Mark").
		Where("mark_names.task_id = ?", taskID).
		Find(&markNames).Error
	if err != nil {
		return nil, nil, err
	}

	for _, markName := range markNames {
		marks = append(marks, markName.Mark)
	}

	return &markNames, &marks, nil
}
