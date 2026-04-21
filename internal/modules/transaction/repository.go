package transaction

import (
	"api-parkir/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	// Fungsi dengan DB Transaction (TX)
	RunInTransaction(fn func(txRepo Repository) error) error
	
	// Operasi Standar
	FindAreaByID(id uint) (*models.AreaParkir, error)
	UpdateArea(area *models.AreaParkir) error
	FindOrCreateKendaraan(kendaraan *models.Kendaraan) error
	FindActiveTransaction(platNomor string) (*models.Transaksi, error)
	CreateTransaction(transaksi *models.Transaksi) error
	UpdateTransaction(transaksi *models.Transaksi) error
	FindTarifByJenis(jenis string) (*models.Tarif, error)
	InsertLogAktivitas(log *models.LogAktivitas) error
	GetDashboardStats() (DashboardStatsResponse, error)
	GetAll() ([]models.Transaksi, error)
	GetLogs() ([]models.LogAktivitas, error)
	GetByDateRange(start, end string) ([]models.Transaksi, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// Wrapper untuk DB Transaction
func (r *repository) RunInTransaction(fn func(txRepo Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)
		return fn(txRepo)
	})
}

func (r *repository) FindAreaByID(id uint) (*models.AreaParkir, error) {
	var area models.AreaParkir
	err := r.db.First(&area, id).Error
	if err != nil {
		return nil, err // FIX: Wajib return nil jika tidak ketemu
	}
	return &area, nil
}

func (r *repository) UpdateArea(area *models.AreaParkir) error {
	return r.db.Save(area).Error
}

func (r *repository) FindOrCreateKendaraan(k *models.Kendaraan) error {
	// Cari berdasarkan plat, kalau tidak ada buat baru
	return r.db.Where(models.Kendaraan{PlatNomor: k.PlatNomor}).
		Assign(models.Kendaraan{JenisKendaraan: k.JenisKendaraan}).
		FirstOrCreate(k).Error
}

func (r *repository) FindActiveTransaction(platNomor string) (*models.Transaksi, error) {
	var trx models.Transaksi
	
	// Kita kembalikan pakai JOIN karena ini yang paling stabil dan sudah terbukti jalan
	err := r.db.Preload("Kendaraan").
		Joins("JOIN tb_kendaraan ON tb_kendaraan.id = tb_transaksi.id_kendaraan").
		Where("tb_kendaraan.plat_nomor = ? AND tb_transaksi.status = ?", platNomor, models.StatusMasuk).
		First(&trx).Error

	if err != nil {
		return nil, err // Wajib return nil agar Golang tidak mengira ada data siluman
	}
	return &trx, nil
}

func (r *repository) CreateTransaction(t *models.Transaksi) error {
	return r.db.Create(t).Error
}

func (r *repository) UpdateTransaction(t *models.Transaksi) error {
	return r.db.Save(t).Error
}

func (r *repository) FindTarifByJenis(jenis string) (*models.Tarif, error) {
	var tarif models.Tarif
	err := r.db.Where("jenis_kendaraan = ?", jenis).First(&tarif).Error
	if err != nil {
		return nil, err // FIX: Wajib return nil jika tidak ketemu
	}
	return &tarif, nil
}

func (r *repository) InsertLogAktivitas(log *models.LogAktivitas) error {
	return r.db.Create(log).Error
}

func (r *repository) GetDashboardStats() (DashboardStatsResponse, error) {
	var stats DashboardStatsResponse

	// A. Hitung kendaraan yang masih parkir (status = 'masuk')
	r.db.Model(&models.Transaksi{}).Where("status = ?", "masuk").Count(&stats.KendaraanParkir)

	// B. Hitung Total Kapasitas dan Jumlah Area Aktif
	type AreaStat struct {
		TotalKapasitas int64
		AreaAktif      int64
	}
	var aStat AreaStat
	r.db.Raw(`
		SELECT COALESCE(SUM(kapasitas), 0) as total_kapasitas, COUNT(id) as area_aktif 
		FROM tb_area_parkir 
		WHERE deleted_at IS NULL
	`).Scan(&aStat)
	stats.TotalKapasitas = aStat.TotalKapasitas
	stats.AreaAktif = aStat.AreaAktif

	// C. Hitung Pendapatan Hari Ini
	r.db.Raw(`
		SELECT COALESCE(SUM(biaya_total), 0) 
		FROM tb_transaksi 
		WHERE DATE(waktu_keluar) = CURRENT_DATE AND deleted_at IS NULL AND status = 'keluar'
	`).Scan(&stats.PendapatanHariIni)

	// ==========================================
	// DATA UNTUK CHART (BARU)
	// ==========================================

	// D. Tren Pendapatan 7 Hari Terakhir
	// Menggunakan TO_CHAR untuk mengubah tanggal jadi format "DD Mon" (misal: "12 Apr")
	r.db.Raw(`
		SELECT 
			TO_CHAR(DATE(waktu_keluar), 'DD Mon') as nama, 
			COALESCE(SUM(biaya_total), 0) as pendapatan 
		FROM tb_transaksi 
		WHERE deleted_at IS NULL 
		AND status = 'keluar'
		AND waktu_keluar >= CURRENT_DATE - INTERVAL '6 days'
		GROUP BY DATE(waktu_keluar)
		ORDER BY DATE(waktu_keluar) ASC
	`).Scan(&stats.TrenPendapatan)

	// Jika kosong (belum ada transaksi 7 hari terakhir), inisialisasi array kosong agar JSON tidak null
	if stats.TrenPendapatan == nil {
		stats.TrenPendapatan = []TrenPendapatan{}
	}

	// E. Utilisasi per Area (Kapasitas vs Terisi)
	// Kita JOIN tb_area_parkir dengan tb_transaksi yang statusnya masih 'masuk'
	r.db.Raw(`
		SELECT 
			a.nama_area as nama, 
			a.kapasitas, 
			COUNT(t.id) as terisi 
		FROM tb_area_parkir a
		LEFT JOIN tb_transaksi t ON a.id = t.id_area AND t.status = 'masuk' AND t.deleted_at IS NULL
		WHERE a.deleted_at IS NULL
		GROUP BY a.id, a.nama_area, a.kapasitas
		ORDER BY a.id ASC
	`).Scan(&stats.UtilisasiArea)

	if stats.UtilisasiArea == nil {
		stats.UtilisasiArea = []UtilisasiArea{}
	}


	return stats, nil
}

func (r *repository) GetAll() ([]models.Transaksi, error) {
	var trxs []models.Transaksi
	// Ambil semua transaksi, urutkan dari yang terbaru ke terlama
	err := r.db.Order("waktu_masuk desc").Find(&trxs).Error
	return trxs, err
}

func (r *repository) GetLogs() ([]models.LogAktivitas, error) {
	var logs []models.LogAktivitas
	// Ambil semua log, urutkan dari yang paling baru
	err := r.db.Order("id desc").Find(&logs).Error
	return logs, err
}

func (r *repository) GetByDateRange(start, end string) ([]models.Transaksi, error) {
    var trxs []models.Transaksi
    // Kita filter berdasarkan waktu_keluar karena rekap biasanya menghitung uang yang sudah diterima (selesai)
    err := r.db.Preload("Kendaraan").
        Where("DATE(waktu_keluar) BETWEEN ? AND ?", start, end).
        Order("waktu_keluar desc").
        Find(&trxs).Error
    return trxs, err
}