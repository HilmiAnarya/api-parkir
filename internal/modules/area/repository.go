package area

import (
	"api-parkir/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(area *models.AreaParkir) error
	FindAll() ([]models.AreaParkir, error)
	FindByID(id uint) (*models.AreaParkir, error)
	Update(area *models.AreaParkir) error
	Delete(area *models.AreaParkir) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db} 
}

func (r *repository) Create(area *models.AreaParkir) error {
	return r.db.Create(area).Error 
}

func (r *repository) FindAll() ([]models.AreaParkir, error)    { 
	var areas []models.AreaParkir
	err := r.db.Find(&areas).Error
	return areas, err
}

func (r *repository) FindByID(id uint) (*models.AreaParkir, error) {
	var area models.AreaParkir
	err := r.db.First(&area, id).Error
	return &area, err
}
func (r *repository) Update(area *models.AreaParkir) error {
	return r.db.Save(area).Error
}

func (r *repository) Delete(area *models.AreaParkir) error {
	return r.db.Delete(area).Error
}