package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB didefinisikan sebagai *gorm.DB agar bisa dipakai di seluruh aplikasi
var DB *gorm.DB

func ConnectDB() {
	// Ambil data dari .env melalui helper GetEnv yang sudah kita buat
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "")
	dbname := GetEnv("DB_NAME", "db_parkir")
	sslmode := GetEnv("DB_SSLMODE", "disable")

	// Merakit DSN secara dinamis
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbname, port, sslmode)

	var err error
	// Langsung buka koneksi menggunakan GORM
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke Database (via GORM): %v", err)
	}

	log.Println("✅ Database terhubung sempurna melalui GORM!")
}