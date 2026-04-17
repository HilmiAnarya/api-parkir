package transaction

import "time"

// Request saat mobil/motor masuk (dari Mesin Kiosk)
type CheckInRequest struct {
	PlatNomor      string `json:"plat_nomor"`
	JenisKendaraan string `json:"jenis_kendaraan"` // "motor" atau "mobil"
	IDArea         uint   `json:"id_area"`
	IDUser         uint   `json:"id_user"` // ID akun petugas Kiosk (sementara ditaruh di body sebelum ada JWT)
	FotoMasuk      string `json:"foto_masuk"`
}

// Request saat keluar (dari Pos Keluar)
type CheckOutRequest struct {
	PlatNomor  string `json:"plat_nomor"`
	IDUser     uint   `json:"id_user"` // ID akun petugas kasir
	FotoKeluar string `json:"foto_keluar"`
}

// Response standar transaksi
type TransactionResponse struct {
	IDParkir    uint       `json:"id_parkir"`
	PlatNomor   string     `json:"plat_nomor"`
	WaktuMasuk  time.Time  `json:"waktu_masuk"`
	WaktuKeluar *time.Time `json:"waktu_keluar,omitempty"`
	DurasiJam   int        `json:"durasi_jam,omitempty"`
	BiayaTotal  float64    `json:"biaya_total,omitempty"`
	Status      string     `json:"status"`
}