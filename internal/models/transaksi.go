package models

import (
	"time"
	"gorm.io/gorm"
)

type StatusTransaksi string

const (
	StatusMasuk  StatusTransaksi = "masuk"
	StatusKeluar StatusTransaksi = "keluar"
)

type Transaksi struct {
	ID           uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	IDKendaraan  uint            `json:"id_kendaraan" gorm:"not null"`
	Kendaraan    Kendaraan       `json:"kendaraan" gorm:"foreignKey:IDKendaraan"`
	WaktuMasuk   time.Time       `json:"waktu_masuk" gorm:"not null"`
	WaktuKeluar  *time.Time      `json:"waktu_keluar"`
	IDTarif      *uint           `json:"id_tarif"`
	Tarif        *Tarif          `json:"tarif" gorm:"foreignKey:IDTarif"`
	DurasiJam    int             `json:"durasi_jam"`
	BiayaTotal   float64         `json:"biaya_total" gorm:"type:decimal(10,0)"`
	Status       StatusTransaksi `json:"status" gorm:"type:status_transaksi;default:'masuk'"`
	IDUser       uint            `json:"id_user" gorm:"not null"`
	User         User            `json:"user" gorm:"foreignKey:IDUser"`
	IDArea       uint            `json:"id_area" gorm:"not null"`
	Area         AreaParkir      `json:"area" gorm:"foreignKey:IDArea"`
	FotoMasuk    string          `json:"foto_masuk" gorm:"size:255;not null"`
	FotoKeluar   string          `json:"foto_keluar" gorm:"size:255"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `json:"-" gorm:"index"`
}

func (Transaksi) TableName() string {
	return "tb_transaksi"
}