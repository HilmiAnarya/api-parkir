package models

import (
	"time"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RolePetugas UserRole = "petugas"
	RoleOwner   UserRole = "owner"
)

type User struct {
	ID           uint           `json:"id_user" gorm:"primaryKey;autoIncrement"`
	NamaLengkap  string         `json:"nama_lengkap" gorm:"size:50;not null"`
	Username     string         `json:"username" gorm:"size:50;uniqueIndex;not null"`
	Password     string         `json:"-" gorm:"size:100;not null"`
	Role         UserRole       `json:"role" gorm:"type:user_role;default:'petugas'"`
	StatusAktif  bool           `json:"status_aktif" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"` // GORM Soft Delete
}

func (User) TableName() string {
	return "tb_user"
}