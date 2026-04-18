package models

import (
	"time"
)

type LogAktivitas struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDUser         uint      `json:"id_user" gorm:"not null"`
	User           User      `json:"user" gorm:"foreignKey:IDUser"`
	Aktivitas      string    `json:"aktivitas" gorm:"size:255;not null"`
	WaktuAktivitas time.Time `json:"waktu_aktivitas" gorm:"default:CURRENT_TIMESTAMP"`
}

func (LogAktivitas) TableName() string {
	return "tb_log_aktivitas"
}