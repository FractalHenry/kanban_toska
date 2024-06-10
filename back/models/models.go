package models

import (
	"time"
)

type User struct {
	Login           string `gorm:"primaryKey; type:varchar(30)"`
	Password        string `gorm:"not null"`
	Email           string `gorm:"not null; unique"`
	EmailVisibility bool   `gorm:"column:email_visibility; default:false; not null"`
	UserDescription string `gorm:"column:user_description; type:varchar(300)"`
}

type RoleOnBoard struct {
	RoleOnBoardID   uint   `gorm:"primaryKey;autoIncrement"`
	RoleOnBoardName string `gorm:"type:varchar(50);not null"`
	SpaceID         uint   `gorm:"not null"`
}

type BoardRoleOnBoard struct {
	RoleOnBoardID uint `gorm:"primaryKey"`
	BoardID       uint `gorm:"primaryKey"`
	CanEdit       bool `gorm:"not null;default:false"`
}

type RoleOnSpace struct {
	RoleOnSpaceID   uint   `gorm:"primaryKey;autoIncrement"`
	RoleOnSpaceName string `gorm:"type:varchar(50);not null"`
	IsAdmin         bool   `gorm:"not null;default:false"`
	IsOwner         bool   `gorm:"not null;default:false"`
	CanEdit         bool   `gorm:"not null;default:false"`
	SpaceID         uint   `gorm:"not null"`
}

type UserRoleOnSpace struct {
	RoleOnSpaceID uint   `gorm:"primaryKey"`
	Login         string `gorm:"primaryKey;type:varchar(30);not null; unique"`
}

type UserBoardRoleOnBoard struct {
	BoardRoleOnBoardID uint   `gorm:"primaryKey"`
	Login              string `gorm:"primaryKey;type:varchar(30);not null; unique"`
}
type Space struct {
	SpaceID        uint   `gorm:"column:space_id;primaryKey;autoIncrement"`
	SpaceName      string `gorm:"column:space_name;type:varchar(70);not null"`
	SpaceInArchive bool   `gorm:"column:space_in_archive;not null;default:false"`
}

type Board struct {
	BoardID        uint   `gorm:"column:board_id;primaryKey;autoIncrement"`
	BoardName      string `gorm:"column:board_name;type:varchar(70);not null"`
	BoardInArchive bool   `gorm:"column:board_in_archive;not null;default:false"`
	SpaceID        uint   `gorm:"column:space_id;not null"`
	Space          Space  `gorm:"foreignKey:SpaceID"`
}

type Card struct {
	CardID        uint   `gorm:"column:card_id;primaryKey;autoIncrement"`
	CardName      string `gorm:"column:card_name;type:varchar(70);not null"`
	CardInArchive bool   `gorm:"column:card_in_archive;not null;default:false"`
	BoardID       uint   `gorm:"column:board_id;not null"`
	Board         Board  `gorm:"foreignKey:BoardID"`
}

type Task struct {
	TaskID        uint   `gorm:"column:task_id;primaryKey;autoIncrement"`
	TaskName      string `gorm:"column:task_name;type:varchar(70);not null"`
	TaskInArchive bool   `gorm:"column:task_in_archive;not null;default:false"`
	CardID        uint   `gorm:"column:card_id;not null"`
	Card          Card   `gorm:"foreignKey:CardID"`
}

type Mark struct {
	MarkID    uint   `gorm:"column:mark_id;primaryKey;autoIncrement"`
	MarkColor string `gorm:"column:mark_color;type:char(7)"`
	TaskID    uint   `gorm:"column:task_id;not null"`
	Task      Task   `gorm:"foreignKey:TaskID"`
}

type MarkName struct {
	MarkID   uint   `gorm:"column:mark_id;not null;uniqueIndex:idx_mark_id"`
	MarkName string `gorm:"column:mark_name;type:varchar(30);not null;primaryKey"`
	Mark     Mark   `gorm:"foreignKey:MarkID"`
}

type Checklist struct {
	ChecklistID   uint   `gorm:"column:checklist_id;primaryKey;autoIncrement"`
	ChecklistName string `gorm:"column:checklist_name;type:varchar(70);not null"`
	TaskID        uint   `gorm:"column:task_id;not null"`
	Task          Task   `gorm:"foreignKey:TaskID"`
}

type ChecklistElement struct {
	ChecklistElementID   uint      `gorm:"column:checklist_element_id;primaryKey;autoIncrement"`
	ChecklistElementName string    `gorm:"column:checklist_element_name;type:varchar(70);not null"`
	Checked              bool      `gorm:"column:checked;not null;default:false"`
	ChecklistID          uint      `gorm:"column:checklist_id;not null"`
	Checklist            Checklist `gorm:"foreignKey:ChecklistID"`
}

type TaskColor struct {
	TaskColor string `gorm:"column:task_color;type:char(7);primaryKey"`
	TaskID    uint   `gorm:"column:task_id;not null"`
	Task      Task   `gorm:"foreignKey:TaskID"`
}

type TaskDescription struct {
	TaskDescription string `gorm:"column:task_description;type:varchar(300);not null;primaryKey"`
	TaskID          uint   `gorm:"column:task_id;not null"`
	Task            Task   `gorm:"foreignKey:TaskID"`
}

type TaskDateStart struct {
	TaskDateStart time.Time `gorm:"column:task_date_start;type:timestamp;not null;primaryKey"`
	TaskID        uint      `gorm:"column:task_id;not null;primaryKey"`
	Task          Task      `gorm:"foreignKey:TaskID"`
}

type TaskNotification struct {
	NotificationID uint          `gorm:"column:notification_id;not null;primaryKey"`
	TaskDateEndID  uint          `gorm:"column:task_date_end_id;not null"`
	TaskDateStart  TaskDateStart `gorm:"foreignKey:TaskDateEndID"`
}

type Notification struct {
	NotificationID          uint   `gorm:"column:notification_id;primaryKey;autoIncrement"`
	NotificationDescription string `gorm:"column:notification_description;type:varchar(50);not null;unique"`
}
