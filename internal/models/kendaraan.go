package models

import (
	"time"
	"gorm.io/gorm"
)

type Kendaraan struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	PlatNomor      string         `json:"plat_nomor" gorm:"size:15;uniqueIndex;not null"`
	JenisKendaraan string         `json:"jenis_kendaraan" gorm:"size:20"`
	Warna          string         `json:"warna" gorm:"size:20"`
	Pemilik        string         `json:"pemilik" gorm:"size:100"`
	IDUser         *uint          `json:"id_user"` // Nullable untuk non-member
	User           *User          `json:"user" gorm:"foreignKey:IDUser"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}