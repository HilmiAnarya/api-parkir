package transaction

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CheckIn(c *fiber.Ctx) error {
	var req CheckInRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	trx, err := h.service.CheckIn(req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Kendaraan berhasil masuk",
		"data": TransactionResponse{
			IDParkir:   trx.ID,
			PlatNomor:  req.PlatNomor,
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