package tarif

import (
	"api-parkir/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(tarif *models.Tarif) error
	FindAll() ([]models.Tarif, error)
	FindByID(id uint) (*models.Tarif, error)
	Update(tarif *models.Tarif) error
	Delete(tarif *models.Tarif) error
}

type repository struct { 
	db *gorm.DB 
}

func NewRepository(db *gorm.DB) Repository { 
	return &repository{db} 
}

func (r *repository) Create(tarif *models.Tarif) error { 
	return r.db.Create(tarif).Error 
}

func (r *repository) FindAll() ([]models.Tarif, error)    { 
	var tarifs []models.Tarif
	err := r.db.Find(&tarifs).Error
	return tarifs, err
}

func (r *repository) FindByID(id uint) (*models.Tarif, error) {
	var tarif models.Tarif
	err := r.db.First(&tarif, id).Error
	return &tarif, err
}

func (r *repository) Update(tarif *models.Tarif) error { 
	return r.db.Save(tarif).Error 
}

func (r *repository) Delete(tarif *models.Tarif) error { 
	return r.db.Delete(tarif).Error 
}