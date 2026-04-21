package user

import (
	"api-parkir/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoute akan merakit repo -> service -> handler secara otomatis
func SetupRoute(router fiber.Router, db *gorm.DB) {
	// Merakit modul User
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	userRoute := router.Group("/users")

	// Rute Publik (Tanpa Token)
	userRoute.Post("/login", handler.Login)
	userRoute.Post("/logout", handler.Logout)

	// Rute CRUD Admin (Wajib Token + Wajib Role Admin)
	adminRoute := userRoute.Group("/", middleware.Protected(), middleware.RequireRole("admin"))
	
	// Kita pindahkan register ke dalam rute admin, sehingga hanya admin yang bisa membuat user baru
	adminRoute.Post("/", handler.CreateUser) 
	adminRoute.Get("/", handler.GetUsers)
	adminRoute.Put("/:id", handler.UpdateUser)
	adminRoute.Delete("/:id", handler.DeleteUser)
}