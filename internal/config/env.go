package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv memuat file .env jika ada.
// Sengaja tidak di-fatal (crash) jika error, karena saat di production (Docker/VPS), 
// file .env mungkin tidak ada dan kita pakai env OS langsung.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: File .env tidak ditemukan, menggunakan Environment Variable dari OS.")
	}
}

// GetEnv mengambil nilai dari environment variable dengan fallback/default value
func GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}