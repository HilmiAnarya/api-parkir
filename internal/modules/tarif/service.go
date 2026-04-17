package tarif

import (
	"api-parkir/internal/models"
	"errors"
)

type Service interface {
	CreateTarif(req TarifRequest) (*models.Tarif, error)
	GetAllTarifs() ([]models.Tarif, error)
	UpdateTarif(id uint, req TarifRequest) (*models.Tarif, error)
	DeleteTarif(id uint) error
}

type service struct { 
	repo Repository 
}

func NewService(repo Repository) Service { 
	return &service{repo} 
}

func (s *service) CreateTarif(req TarifRequest) (*models.Tarif, error) {
	tarif := models.Tarif{ 
		JenisKendaraan: models.JenisKendaraan(req.JenisKendaraan), 
		TarifPerJam: req.TarifPerJam,
	}
	err := s.repo.Create(&tarif)
	return &tarif, err
}

func (s *service) GetAllTarifs() ([]models.Tarif, error) { 
	return s.repo.FindAll() 
}

func (s *service) UpdateTarif(id uint, req TarifRequest) (*models.Tarif, error) {
	tarif, err := s.repo.FindByID(id)
	if err != nil { return nil, errors.New("tarif tidak ditemukan") }

	tarif.JenisKendaraan = models.JenisKendaraan(req.JenisKendaraan)
	tarif.TarifPerJam = req.TarifPerJam
	err = s.repo.Update(tarif)
	return tarif, err
}

func (s *service) DeleteTarif(id uint) error {
	tarif, err := s.repo.FindByID(id)
	if err != nil { return errors.New("tarif tidak ditemukan") }
	return s.repo.Delete(tarif)
}