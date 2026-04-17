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
	err := r.db.Preload("Kendaraan").
		Joins("JOIN tb_kendaraan ON tb_kendaraan.id = tb_transaksi.id_kendaraan").
		Where("tb_kendaraan.plat_nomor = ? AND tb_transaksi.status = ?", platNomor, models.StatusMasuk).
		First(&trx).Error
	
	if err != nil {
		return nil, err // FIX: Mencegah service mengira data ada (menghindari 409 Conflict palsu)
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