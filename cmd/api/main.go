package main

import (
	"api-parkir/internal/config"
	"api-parkir/internal/modules/area"
	"api-parkir/internal/modules/tarif"
	"api-parkir/internal/modules/transaction"
	"api-parkir/internal/modules/user"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Gagal memuat timezone: %v", err)
	}
	time.Local = loc
	// 1. Load Environment Variables
	config.LoadEnv()

	// 2. Inisialisasi Database (Sekarang sudah GORM & Env-based)
	config.ConnectDB()

	// 3. Setup Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sistem Parkir API v1.0",
	})

	// ==========================================
	// ⚡ TAMBAHKAN MIDDLEWARE CORS DI SINI ⚡
	// ==========================================
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Mengizinkan semua origin (frontend mana pun)
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	app.Use(logger.New())

	// 4. Routing
	api := app.Group("/api")

	// Kirim config.DB ke modul-modul yang membutuhkan
	user.SetupRoute(api, config.DB)
	// Daftarkan Module Transaction
	transaction.SetupRoute(api, config.DB)
	// Daftarkan Module Area
	area.SetupRoute(api, config.DB)
	// Daftarkan Module Tarif
	tarif.SetupRoute(api, config.DB)

	// 5. Start Server
	port := config.GetEnv("PORT", "8080")
	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Server crash: %v", err)
	}
}