package transaction

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"api-parkir/internal/utils" // Import folder utils yang baru dibuat
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CheckIn(c *fiber.Ctx) error {
	// 1. Tangkap File Foto dari Request Multipart
	file, err := c.FormFile("foto")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Foto wajib diunggah!"})
	}

	// Tangkap data teks lainnya dari form
	jenisKendaraan := c.FormValue("jenis_kendaraan")
	idArea, _ := strconv.Atoi(c.FormValue("id_area"))
	idUser, _ := strconv.Atoi(c.FormValue("id_user"))

	if jenisKendaraan == "" || idArea == 0 || idUser == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data jenis_kendaraan, id_area, dan id_user wajib diisi!"})
	}

	// 2. Buka file untuk dibaca menjadi bytes
	fileHeader, _ := file.Open()
	defer fileHeader.Close()
	fileBytes, _ := io.ReadAll(fileHeader)

	// ========================================================
	// 3. TUGAS AI: Kirim foto ke Python dan dapatkan Plat Nomor
	// ========================================================
	platNomorAI, err := utils.DetectPlate(fileBytes, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}

	// 4. Simpan foto ke folder lokal (Sebagai Bukti CCTV)
	// Buat folder "uploads" otomatis jika belum ada
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", os.ModePerm)
	}
	
	// Format nama file: UNIXTIME_PLATNOMOR.jpg
	savePath := fmt.Sprintf("./uploads/%d_%s_%s", time.Now().Unix(), platNomorAI, file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan file gambar lokal"})
	}

	// 5. Rakit Request untuk dikirim ke Service Database
	req := CheckInRequest{
		PlatNomor:      platNomorAI, // Hasil murni dari Python AI!
		JenisKendaraan: jenisKendaraan,
		IDArea:         uint(idArea),
		IDUser:         uint(idUser),
		FotoMasuk:      savePath,
	}

	// 6. Eksekusi Jantung Transaksi (Logika Database)
	trx, err := h.service.CheckIn(req)
	if err != nil {
		// Jika database menolak (misal: area penuh/double check-in), kita bisa opsional menghapus foto agar hardisk tidak penuh
		_ = os.Remove(savePath)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	// 7. Berhasil! Kembalikan ke Frontend
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Kendaraan berhasil masuk via AI Detection",
		"data": TransactionResponse{
			IDParkir:   trx.ID,
			PlatNomor:  req.PlatNomor, // Akan terlihat "B1234XYZ" atau sesuai foto
			WaktuMasuk: trx.WaktuMasuk,
			Status:     string(trx.Status),
		},
	})
}

func (h *Handler) CheckOut(c *fiber.Ctx) error {
	var req CheckOutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	trx, err := h.service.CheckOut(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Transaksi keluar berhasil",
		"data": TransactionResponse{
			IDParkir:    trx.ID,
			PlatNomor:   req.PlatNomor,
			WaktuMasuk:  trx.WaktuMasuk,
			WaktuKeluar: trx.WaktuKeluar,
			DurasiJam:   trx.DurasiJam,
			BiayaTotal:  trx.BiayaTotal,
			Status:      string(trx.Status),
		},
	})
}

func (h *Handler) CheckPrice(c *fiber.Ctx) error {
	platNomor := c.Params("plat_nomor")
	if platNomor == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Plat nomor wajib diisi"})
	}

	trx, err := h.service.CheckPrice(platNomor)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Tagihan berhasil dihitung",
		"data": TransactionResponse{
			IDParkir:   trx.ID,
			PlatNomor:  platNomor,
			WaktuMasuk: trx.WaktuMasuk,
			DurasiJam:  trx.DurasiJam,
			BiayaTotal: trx.BiayaTotal,
		},
	})
}

func (h *Handler) GetDashboardStats(c *fiber.Ctx) error {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data statistik"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	trxs, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil laporan transaksi"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    trxs,
	})
}

func (h *Handler) GetLogs(c *fiber.Ctx) error {
	logs, err := h.service.GetLogs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil log aktivitas"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    logs,
	})
}

func (h *Handler) ExportExcel(c *fiber.Ctx) error {
    start := c.Query("start_date")
    end := c.Query("end_date")

    if start == "" || end == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Rentang tanggal harus diisi"})
    }

    trxs, err := h.service.GetByDateRange(start, end)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    f := excelize.NewFile()
    defer f.Close()
    sheet := "Rekap Parkir"
    f.SetSheetName("Sheet1", sheet)

    // Header
    headers := []string{"ID", "Plat Nomor", "Waktu Masuk", "Waktu Keluar", "Durasi (Jam)", "Total Biaya"}
    for i, h := range headers {
        cell, _ := excelize.CoordinatesToCellName(i+1, 1)
        f.SetCellValue(sheet, cell, h)
    }

    // Data Rows
    for i, t := range trxs {
        row := i + 2
        f.SetCellValue(sheet, fmt.Sprintf("A%d", row), t.ID)
        f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t.Kendaraan.PlatNomor)
        f.SetCellValue(sheet, fmt.Sprintf("C%d", row), t.WaktuMasuk.Format("2006-01-02 15:04"))
        if t.WaktuKeluar != nil {
            f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.WaktuKeluar.Format("2006-01-02 15:04"))
        }
        f.SetCellValue(sheet, fmt.Sprintf("E%d", row), t.DurasiJam)
        f.SetCellValue(sheet, fmt.Sprintf("F%d", row), t.BiayaTotal)
    }

    c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=rekap_parkir_%s_%s.xlsx", start, end))

    return f.Write(c.Response().BodyWriter())
}