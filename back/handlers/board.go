package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetBoardDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID доски из пути запроса
	vars := mux.Vars(r)
	boardID, err := strconv.ParseUint(vars["boardId"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID доски", http.StatusBadRequest)
		return
	}

	// Получаем детальную информацию о доске
	board, cards, _, infoBlocks, err := repo.GetBoardDetails(uint(boardID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем пользователей и их роли на доске
	boardUsers, err := repo.GetBoardUsers(uint(boardID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	space, _ := repo.GetSpaceByBoard(uint(boardID))
	spaceUsers, err := repo.GetSpaceUsers(space.SpaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type MarkNameWithMark struct {
		MarkName models.MarkName `json:"markName"`
		Mark     models.Mark     `json:"mark"`
	}

	type ChecklistWithElements struct {
		Checklist         models.Checklist          `json:"checklist"`
		ChecklistElements []models.ChecklistElement `json:"checklistElements"`
	}

	type TaskWithDetails struct {
		Task              models.Task              `json:"task"`
		TaskColor         string                   `json:"taskColor,omitempty"`
		TaskDescription   string                   `json:"taskDescription,omitempty"`
		TaskDateStart     *models.TaskDateStart    `json:"taskDateStart,omitempty"`
		TaskDateEnd       *models.TaskDateEnd      `json:"taskDateEnd,omitempty"`
		TaskNotifications *[]models.Notification   `json:"taskNotifications,omitempty"`
		TaskMarkNames     *[]MarkNameWithMark      `json:"taskMarkNames,omitempty"`
		Checklists        *[]ChecklistWithElements `json:"checklists,omitempty"`
	}

	type CardWithTasks struct {
		Card  models.Card       `json:"card"`
		Tasks []TaskWithDetails `json:"tasks"`
	}

	// Создаем ответ в формате JSON
	response := struct {
		ID         uint                        `json:"id"`
		Name       string                      `json:"name"`
		Cards      []CardWithTasks             `json:"cards"`
		InfoBlocks []models.InformationalBlock `json:"infoBlocks"`
		BoardUsers []struct {
			Login   string `json:"login"`
			CanEdit bool   `json:"can_edit"`
		} `json:"boardUsers"`
		SpaceUsers []struct {
			Login   string `json:"login"`
			IsAdmin bool   `json:"is_admin"`
			IsOwner bool   `json:"is_owner"`
			CanEdit bool   `json:"can_edit"`
		} `json:"spaceUsers"`
	}{
		ID:         board.BoardID,
		Name:       board.BoardName,
		Cards:      make([]CardWithTasks, 0),
		InfoBlocks: infoBlocks,
		BoardUsers: make([]struct {
			Login   string `json:"login"`
			CanEdit bool   `json:"can_edit"`
		}, 0),
		SpaceUsers: make([]struct {
			Login   string `json:"login"`
			IsAdmin bool   `json:"is_admin"`
			IsOwner bool   `json:"is_owner"`
			CanEdit bool   `json:"can_edit"`
		}, 0),
	}

	// Заполняем информацию о карточках и задачах
	for _, card := range cards {
		var tasksWithDetails []TaskWithDetails
		cardTasks, err := repo.GetCardTasks(card.CardID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, task := range cardTasks {
			taskColor, _ := repo.GetTaskColor(task.TaskID)
			taskDescription, _ := repo.GetTaskDescription(task.TaskID)
			taskDateStart, _ := repo.GetTaskStartDate(task.TaskID)
			taskDateEnd, _ := repo.GetTaskEndDate(task.TaskID)
			taskNotifications, _ := repo.GetTaskNotifications(task.TaskID)
			taskMarkNames, _ := repo.GetTaskMarkNames(task.TaskID)
			var markNamesWithMarks []MarkNameWithMark
			if taskMarkNames != nil {
				for _, markName := range *taskMarkNames {
					markNamesWithMarks = append(markNamesWithMarks, MarkNameWithMark{
						MarkName: markName,
						Mark:     markName.Mark,
					})
				}
			}
			checklists, _ := repo.GetChecklistsByTaskID(task.TaskID)
			checklistWithElements := []ChecklistWithElements{}
			if checklists != nil {
				for _, checklist := range *checklists {
					checklistElements, err := repo.GetChecklistElementsByChecklistID(checklist.ChecklistID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					checklistWithElements = append(checklistWithElements, ChecklistWithElements{
						Checklist:         checklist,
						ChecklistElements: *checklistElements,
					})
				}
			}

			tasksWithDetails = append(tasksWithDetails, TaskWithDetails{
				Task:              task,
				TaskColor:         taskColor,
				TaskDescription:   taskDescription,
				TaskDateStart:     taskDateStart,
				TaskDateEnd:       taskDateEnd,
				TaskNotifications: taskNotifications,
				TaskMarkNames:     &markNamesWithMarks,
				Checklists:        &checklistWithElements,
			})
		}

		response.Cards = append(response.Cards, CardWithTasks{
			Card:  card,
			Tasks: tasksWithDetails,
		})
	}

	// Заполняем информацию о пользователях
	if boardUsers != nil {
		for _, user := range *boardUsers {
			response.BoardUsers = append(response.BoardUsers, struct {
				Login   string `json:"login"`
				CanEdit bool   `json:"can_edit"`
			}{
				Login:   user.Login,
				CanEdit: user.RoleOnBoard.CanEdit,
			})
		}
	}

	if spaceUsers != nil {
		for _, user := range *spaceUsers {
			response.SpaceUsers = append(response.SpaceUsers, struct {
				Login   string `json:"login"`
				IsAdmin bool   `json:"is_admin"`
				IsOwner bool   `json:"is_owner"`
				CanEdit bool   `json:"can_edit"`
			}{
				Login:   user.Login,
				IsAdmin: user.RoleOnSpace.IsAdmin,
				IsOwner: user.RoleOnSpace.IsOwner,
				CanEdit: user.RoleOnSpace.CanEdit,
			})
		}
	}

	// Отправляем ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Header.Get("login")

	var boardData struct {
		SpaceID   uint   `json:"spaceid"`
		BoardName string `json:"boardname"`
	}

	err := json.NewDecoder(r.Body).Decode(&boardData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	board := &models.Board{
		SpaceID:   boardData.SpaceID,
		BoardName: boardData.BoardName,
	}

	err = repo.CreateBoard(board, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
