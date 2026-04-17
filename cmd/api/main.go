package main

import (
	"log"
	"api-parkir/internal/config" 
	"api-parkir/internal/modules/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Load Environment Variables
	config.LoadEnv()

	// 2. Inisialisasi Database (Sekarang sudah GORM & Env-based)
	config.ConnectDB()

	// 3. Setup Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sistem Parkir API v1.0",
	})

	app.Use(logger.New())

	// 4. Routing
	api := app.Group("/api")

	// Kirim config.DB ke modul-modul yang membutuhkan
	user.SetupRoute(api, config.DB)

	// 5. Start Server
	port := config.GetEnv("PORT", "8080")
	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Server crash: %v", err)
	}
}