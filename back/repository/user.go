package repository

import (
	"backend/models"

	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return r.db.Create(user).Error
}

func (r *Repository) FindUserByLogin(login string) (*models.User, error) {
	var user models.User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateEmailVisibility(login string, emailVisibility bool) error {
	return r.db.Model(&models.User{}).Where("login = ?", login).Update("email_visibility", emailVisibility).Error
}

func (r *Repository) UpdateUserDescription(login string, description string) error {
	return r.db.Model(&models.User{}).Where("login = ?", login).Update("user_description", description).Error
}
