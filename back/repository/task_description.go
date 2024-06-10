package repository

import (
	"backend/models"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

func (r *Repository) UpdateTaskColor(taskID uint, newColor string, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var task models.Task
	err = r.db.First(&task, taskID).Error
	if err != nil {
		return err
	}

	var card models.Card
	err = r.db.First(&card, task.CardID).Error
	if err != nil {
		return err
	}

	var board models.Board
	err = r.db.First(&board, card.BoardID).Error
	if err != nil {
		return err
	}

	var currentRoleOnSpace models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&currentRoleOnSpace).Error
	if err != nil {
		return err
	}

	var currentBoardRole models.BoardRoleOnBoard
	err = r.db.Model(&models.BoardRoleOnBoard{}).
		Where("board_id = ? AND role_on_board_id IN (SELECT role_on_board_id FROM user_board_role_on_boards WHERE login = ?)", board.BoardID, user.Login).
		First(&currentBoardRole).Error
	if err != nil {
		return err
	}

	if currentRoleOnSpace.IsOwner || currentRoleOnSpace.IsAdmin || currentRoleOnSpace.CanEdit || currentBoardRole.CanEdit {
		taskColor := models.TaskColor{
			TaskColor: newColor,
			TaskID:    taskID,
		}
		return r.db.Save(&taskColor).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для изменения цвета задания")
	}
}

func (r *Repository) DeleteTaskDescription(taskID uint, userLogin string) error {
	user, err := r.FindUserByLogin(userLogin)
	if err != nil {
		return err
	}

	var task models.Task
	err = r.db.First(&task, taskID).Error
	if err != nil {
		return err
	}

	var card models.Card
	err = r.db.First(&card, task.CardID).Error
	if err != nil {
		return err
	}

	var board models.Board
	err = r.db.First(&board, card.BoardID).Error
	if err != nil {
		return err
	}

	var currentRoleOnSpace models.RoleOnSpace
	err = r.db.Model(&models.RoleOnSpace{}).
		Where("space_id = ? AND role_on_space_id IN (SELECT role_on_space_id FROM user_role_on_spaces WHERE login = ?)", board.SpaceID, user.Login).
		First(&currentRoleOnSpace).Error
	if err != nil {
		return err
	}

	var currentBoardRole models.BoardRoleOnBoard
	err = r.db.Model(&models.BoardRoleOnBoard{}).
		Where("board_id = ? AND role_on_board_id IN (SELECT role_on_board_id FROM user_board_role_on_boards WHERE login = ?)", board.BoardID, user.Login).
		First(&currentBoardRole).Error
	if err != nil {
		return err
	}

	if currentRoleOnSpace.IsOwner || currentRoleOnSpace.IsAdmin || currentRoleOnSpace.CanEdit || currentBoardRole.CanEdit {
		taskDescription := models.TaskDescription{
			TaskID: taskID,
		}
		return r.db.Delete(&taskDescription).Error
	} else {
		return fmt.Errorf("у пользователя нет прав для удаления описания задания")
	}
}

func (r *Repository) GetTaskDescription(taskID uint) (string, error) {
	var taskDescription models.TaskDescription
	err := r.db.Where("task_id = ?", taskID).First(&taskDescription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // Если описания не существует, возвращаем пустую строку
		}
		return "", err
	}
	return taskDescription.TaskDescription, nil
}
