package models

import (
	"time"
	"gorm.io/gorm"
)

type AreaParkir struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaArea  string         `json:"nama_area" gorm:"size:50;not null"`
	Kapasitas int            `json:"kapasitas" gorm:"default:0"`
	Terisi    int            `json:"terisi" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}