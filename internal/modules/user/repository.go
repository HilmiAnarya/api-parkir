package user

import (
	"api-parkir/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

// Constructor untuk membuat instance repository baru
func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *repository) Delete(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}