package area

import (
	"api-parkir/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoute(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	areaRoute := router.Group("/area")

	areaRoute.Use(middleware.Protected())
	areaRoute.Use(middleware.RequireRole("admin"))

	areaRoute.Post("/", handler.Create)
	areaRoute.Get("/", handler.GetAll)
	areaRoute.Put("/:id", handler.Update)
	areaRoute.Delete("/:id", handler.Delete)
}