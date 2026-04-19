package transaction

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoute(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// Group route: http://localhost:8080/api/transactions
	trxRoute := router.Group("/transactions")

	trxRoute.Get("/all", handler.GetAll)
	trxRoute.Post("/in", handler.CheckIn)
	trxRoute.Post("/out", handler.CheckOut)
	trxRoute.Get("/price/:plat_nomor", handler.CheckPrice)
	trxRoute.Get("/stats/dashboard", handler.GetDashboardStats)
	trxRoute.Get("/logs", handler.GetLogs)
}