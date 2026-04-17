package models

import (
	"time"
	"gorm.io/gorm"
)

type JenisKendaraan string

const (
	Motor   JenisKendaraan = "motor"
	Mobil   JenisKendaraan = "mobil"
	Lainnya JenisKendaraan = "lainnya"
)

type Tarif struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	JenisKendaraan JenisKendaraan `json:"jenis_kendaraan" gorm:"type:jenis_kendaraan_enum;not null"`
	TarifPerJam    float64        `json:"tarif_per_jam" gorm:"type:decimal(10,0);not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}