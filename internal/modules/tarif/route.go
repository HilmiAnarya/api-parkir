package tarif

import (
	"api-parkir/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoute(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	tarifRoute := router.Group("/tarif")
	
	tarifRoute.Use(middleware.Protected())
	tarifRoute.Use(middleware.RequireRole("admin"))

	tarifRoute.Post("/", handler.Create)
	tarifRoute.Get("/", handler.GetAll)
	tarifRoute.Put("/:id", handler.Update)
	tarifRoute.Delete("/:id", handler.Delete)
}