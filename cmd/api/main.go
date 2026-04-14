package main

import (
	"log"
	"api-parkir/internal/config" // Sesuaikan "api-parkir" dengan nama module di go.mod kamu

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Load Environment Variables
	config.LoadEnv()

	// 2. Hubungkan ke Database PostgreSQL
	config.ConnectDB()

	// 3. Inisialisasi Go Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sistem Parkir API v1.0",
	})

	// Tambahkan middleware logger agar setiap request HTTP terlihat di terminal
	app.Use(logger.New())

	// Endpoint sederhana untuk mengecek server hidup atau tidak
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
			"pesan":  "Server Go Fiber dan Database siap menerima perintah!",
		})
	})

	// 4. Jalankan Server
	port := config.GetEnv("PORT", "8080")
	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("❌ Server gagal berjalan: %v", err)
	}
}