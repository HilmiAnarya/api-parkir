package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoute akan merakit repo -> service -> handler secara otomatis
func SetupRoute(router fiber.Router, db *gorm.DB) {
	// Merakit modul User
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// Membuat grup route khusus user
	// Contoh jadinya: http://localhost:8080/api/users
	userRoute := router.Group("/users")

	// Endpoint
	userRoute.Post("/register", handler.CreateUser) // Menambah user (Admin only nanti)
	userRoute.Post("/login", handler.Login)         // Login sistem
	userRoute.Get("/", handler.GetUsers)            // Menampilkan semua user
}