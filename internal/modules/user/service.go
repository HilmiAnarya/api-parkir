package user

import (
	"api-parkir/internal/config"
	"api-parkir/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(req CreateUserRequest) (*models.User, error)
	Login(req LoginRequest) (string, *models.User, error)
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

func (s *service) Login(req LoginRequest) (string, *models.User, error) {
	// 1. Cari user berdasarkan username
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return "", nil, errors.New("username atau password salah")
	}

	// 2. VERIFIKASI BCRYPT (Perbaikan Utama)
	// Kita bandingkan hash di DB dengan password murni dari request
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Jika error, berarti password salah
		return "", nil, errors.New("username atau password salah")
	}

	// 3. Buat Token JWT (Jika password cocok)
	claims := jwt.MapClaims{
		"id_user": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := config.GetEnv("JWT_SECRET", "rahasia_parkir_2026")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, errors.New("gagal menerbitkan token")
	}

	return tokenString, user, nil
}

func (s *service) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}