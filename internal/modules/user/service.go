package user

import (
	"api-parkir/internal/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(req CreateUserRequest) (*models.User, error)
	Login(req LoginRequest) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) RegisterUser(req CreateUserRequest) (*models.User, error) {
	// Cek apakah username sudah dipakai
	existingUser, _ := s.repo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username sudah digunakan")
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	// Rakit model user
	newUser := models.User{
		NamaLengkap: req.NamaLengkap,
		Username:    req.Username,
		Password:    string(hashedPassword),
		Role:        models.UserRole(req.Role),
		StatusAktif: true,
	}

	// Simpan ke DB
	err = s.repo.Create(&newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *service) Login(req LoginRequest) (*models.User, error) {
	// Cari user berdasarkan username
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	if !user.StatusAktif {
		return nil, errors.New("akun anda sedang dinonaktifkan")
	}

	// Cocokkan password plain dari request dengan password hash dari DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	return user, nil
}

func (s *service) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}