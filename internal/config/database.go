package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Wajib ada underscore (init) agar driver PostgreSQL terbaca
)

// DB adalah variabel global agar bisa diakses oleh Repository nanti
var DB *sql.DB

// ConnectDB membuka koneksi ke PostgreSQL
func ConnectDB() {
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "")
	dbname := GetEnv("DB_NAME", "db_parkir")
	sslmode := GetEnv("DB_SSLMODE", "disable")

	// Merakit Connection String
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Gagal membuka database: %v\n", err)
	}

	// Tes koneksi sebenarnya (Ping)
	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Gagal terhubung (Ping) ke PostgreSQL: %v\n", err)
	}

	// Connection Pooling (Sangat penting agar server tidak crash saat banyak mobil masuk bersamaan)
	DB.SetMaxOpenConns(25)                 // Maksimal 25 koneksi aktif bersamaan
	DB.SetMaxIdleConns(25)                 // Maksimal 25 koneksi standby
	DB.SetConnMaxLifetime(5 * time.Minute) // Tiap koneksi direfresh setiap 5 menit

	log.Println("✅ Berhasil terhubung ke PostgreSQL!")
}