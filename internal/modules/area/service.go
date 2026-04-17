package area

import (
	"api-parkir/internal/models"
	"errors"
)

type Service interface {
	CreateArea(req AreaRequest) (*models.AreaParkir, error)
	GetAllAreas() ([]models.AreaParkir, error)
	UpdateArea(id uint, req AreaRequest) (*models.AreaParkir, error)
	DeleteArea(id uint) error
}

type service struct { repo Repository }

func NewService(repo Repository) Service {
	return &service{repo} 
}

func (s *service) CreateArea(req AreaRequest) (*models.AreaParkir, error) {
	area := models.AreaParkir{ NamaArea: req.NamaArea, Kapasitas: req.Kapasitas, Terisi: 0 }
	err := s.repo.Create(&area)
	return &area, err
}

func (s *service) GetAllAreas() ([]models.AreaParkir, error) {
	return s.repo.FindAll() 
}

func (s *service) UpdateArea(id uint, req AreaRequest) (*models.AreaParkir, error) {
	area, err := s.repo.FindByID(id)
	if err != nil { return nil, errors.New("area tidak ditemukan") }

	area.NamaArea = req.NamaArea
	area.Kapasitas = req.Kapasitas
	err = s.repo.Update(area)
	return area, err
}

func (s *service) DeleteArea(id uint) error {
	area, err := s.repo.FindByID(id)
	if err != nil { return errors.New("area tidak ditemukan") }
	return s.repo.Delete(area)
}