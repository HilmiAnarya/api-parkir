package transaction

import (
	"api-parkir/internal/models"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	CheckIn(req CheckInRequest) (*models.Transaksi, error)
	CheckOut(req CheckOutRequest) (*models.Transaksi, error)
	CheckPrice(platNomor string) (*models.Transaksi, error)
	GetDashboardStats() (DashboardStatsResponse, error)
	GetAll() ([]models.Transaksi, error)
	GetLogs() ([]models.LogAktivitas, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CheckIn(req CheckInRequest) (*models.Transaksi, error) {
	var newTrx models.Transaksi

	// Gunakan Database Transaction agar aman (ACID)
	err := s.repo.RunInTransaction(func(txRepo Repository) error {
		// 1. Cek Kapasitas Area
		area, err := txRepo.FindAreaByID(req.IDArea)
		if err != nil {
			return errors.New("area parkir tidak ditemukan")
		}
		if area.Terisi >= area.Kapasitas {
			return errors.New("mohon maaf, area parkir penuh")
		}

		// 2. Cek apakah kendaraan sudah parkir tapi belum keluar
		activeTrx, err := txRepo.FindActiveTransaction(req.PlatNomor)
		
		// PENJAGA KEAMANAN UTAMA:
		// Jika ada error, DAN errornya bukan "Data Tidak Ditemukan", hentikan transaksi!
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("terjadi kesalahan database saat mengecek status kendaraan") 
		}
		
		// Jika mobilnya memang ada di dalam (activeTrx tidak kosong)
		if activeTrx != nil {
			return errors.New("kendaraan ini sedang berada di dalam area parkir")
		}

		// 3. Daftarkan/Ambil Data Kendaraan
		kendaraan := models.Kendaraan{
			PlatNomor:      req.PlatNomor,
			JenisKendaraan: req.JenisKendaraan,
		}
		if err := txRepo.FindOrCreateKendaraan(&kendaraan); err != nil {
			return errors.New("gagal memproses data kendaraan")
		}

		// 4. Buat Transaksi Baru
		newTrx = models.Transaksi{
			IDKendaraan: kendaraan.ID,
			WaktuMasuk:  time.Now(),
			Status:      models.StatusMasuk,
			IDUser:      req.IDUser,
			IDArea:      req.IDArea,
			FotoMasuk:   req.FotoMasuk,
		}
		if err := txRepo.CreateTransaction(&newTrx); err != nil {
			return errors.New("gagal mencetak tiket masuk")
		}

		// 5. Tambah Kapasitas Terisi
		area.Terisi += 1
		if err := txRepo.UpdateArea(area); err != nil {
			return errors.New("gagal mengupdate kapasitas area")
		}

		// ==========================================
		// 6. CATAT LOG AKTIVITAS (BARU)
		// ==========================================
		logMasuk := models.LogAktivitas{
			IDUser:    req.IDUser,
			Aktivitas: "Memproses Check-In kendaraan plat: " + req.PlatNomor,
		}
		// Abaikan error log agar transaksi utama tidak gagal jika log gagal dicatat
		_ = txRepo.InsertLogAktivitas(&logMasuk)

		return nil // Sukses
	})

	if err != nil {
		return nil, err
	}
	return &newTrx, nil
}

func (s *service) CheckOut(req CheckOutRequest) (*models.Transaksi, error) {
	var trxToUpdate *models.Transaksi

	err := s.repo.RunInTransaction(func(txRepo Repository) error {
		// 1. Cari transaksi yang masih gantung (belum keluar)
		trx, err := txRepo.FindActiveTransaction(req.PlatNomor)
		if err != nil {
			return errors.New("tiket masuk tidak ditemukan atau kendaraan sudah keluar")
		}

		// 2. Ambil Tarif berdasarkan jenis kendaraan (motor/mobil)
		tarif, err := txRepo.FindTarifByJenis(trx.Kendaraan.JenisKendaraan)
		if err != nil {
			return errors.New("master tarif belum diatur untuk jenis kendaraan ini")
		}

		// 3. Hitung Durasi dan Biaya
		waktuKeluar := time.Now()
		durasiAsli := waktuKeluar.Sub(trx.WaktuMasuk).Hours()
		
		// Pembulatan ke atas (misal 1.2 jam = 2 jam)
		durasiJam := int(math.Ceil(durasiAsli))
		if durasiJam < 1 {
			durasiJam = 1 // Minimal bayar 1 jam
		}
		biayaTotal := float64(durasiJam) * tarif.TarifPerJam

		// 4. Update Transaksi
		trx.WaktuKeluar = &waktuKeluar
		trx.IDTarif = &tarif.ID
		trx.DurasiJam = durasiJam
		trx.BiayaTotal = biayaTotal
		trx.Status = models.StatusKeluar
		trx.FotoKeluar = req.FotoKeluar

		if err := txRepo.UpdateTransaction(trx); err != nil {
			return errors.New("gagal menyimpan transaksi keluar")
		}

		// 5. Kurangi Kapasitas Terisi di Area
		area, _ := txRepo.FindAreaByID(trx.IDArea)
		if area != nil && area.Terisi > 0 {
			area.Terisi -= 1
			txRepo.UpdateArea(area)
		}

		// ==========================================
		// 6. CATAT LOG AKTIVITAS (BARU)
		// ==========================================
		logKeluar := models.LogAktivitas{
			IDUser:    req.IDUser,
			Aktivitas: "Memproses Check-Out kendaraan plat: " + req.PlatNomor,
		}
		_ = txRepo.InsertLogAktivitas(&logKeluar)

		trxToUpdate = trx
		return nil
	})

	if err != nil {
		return nil, err
	}
	return trxToUpdate, nil
}

func (s *service) CheckPrice(platNomor string) (*models.Transaksi, error) {
	// 1. Cari transaksi yang masih gantung (belum keluar)
	trx, err := s.repo.FindActiveTransaction(platNomor)
	if err != nil {
		return nil, errors.New("tiket masuk tidak ditemukan atau kendaraan sudah keluar")
	}

	// 2. Ambil Tarif berdasarkan jenis kendaraan (motor/mobil)
	tarif, err := s.repo.FindTarifByJenis(trx.Kendaraan.JenisKendaraan)
	if err != nil {
		return nil, errors.New("master tarif belum diatur")
	}

	// 3. Hitung Durasi dan Biaya "Saat Ini" (Virtual)
	waktuSekarang := time.Now()
	durasiAsli := waktuSekarang.Sub(trx.WaktuMasuk).Hours()
	
	durasiJam := int(math.Ceil(durasiAsli))
	if durasiJam < 1 {
		durasiJam = 1
	}
	biayaTotal := float64(durasiJam) * tarif.TarifPerJam

	// Tempelkan hasil hitungan sementara ke struct (TIDAK DISIMPAN KE DB)
	trx.DurasiJam = durasiJam
	trx.BiayaTotal = biayaTotal

	return trx, nil
}

func (s *service) GetDashboardStats() (DashboardStatsResponse, error) {
	return s.repo.GetDashboardStats()
}

func (s *service) GetAll() ([]models.Transaksi, error) {
	return s.repo.GetAll()
}

func (s *service) GetLogs() ([]models.LogAktivitas, error) {
	return s.repo.GetLogs()
}