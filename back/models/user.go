package models

type User struct {
	ID       uint   `gorm:"primaryKey; autoIncrement"`
	Login    string `gorm:"not null; unique"`
	Password string `gorm:"not null"`
	Email    string `gorm:"not null; unique"`
}
