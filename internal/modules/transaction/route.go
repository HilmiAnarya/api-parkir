package transaction

import (
	"api-parkir/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoute(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// Group route: http://localhost:8080/api/transactions
	trxRoute := router.Group("/transactions")

	// Rute Laporan bisa diakses Petugas DAN Admin
	trxRoute.Get("/all", middleware.Protected(), middleware.RequireRole("admin", "petugas", "owner"), handler.GetAll)

	// Rute Dashboard Stats dan Logs HANYA Admin
	trxRoute.Get("/stats/dashboard", middleware.Protected(), middleware.RequireRole("admin", "owner"), handler.GetDashboardStats)
	trxRoute.Get("/export/excel", middleware.Protected(), middleware.RequireRole("admin", "owner"), handler.ExportExcel)
	trxRoute.Get("/logs", middleware.Protected(), middleware.RequireRole("admin"), handler.GetLogs)
	
	// Rute In/Out/Price biarkan publik dulu untuk kemudahan tes Kiosk (Atau pasang Protected jika Kiosk nanti diset login)
	trxRoute.Post("/in", handler.CheckIn)
	trxRoute.Post("/out", handler.CheckOut)
	trxRoute.Get("/price/:plat_nomor", handler.CheckPrice)
}