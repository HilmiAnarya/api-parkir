package user

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	user, err := h.service.RegisterUser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil membuat user baru",
		"data": UserResponse{
			ID:          user.ID,
			NamaLengkap: user.NamaLengkap,
			Username:    user.Username,
			Role:        string(user.Role),
			StatusAktif: user.StatusAktif,
		},
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	// Panggil service Login
	token, user, err := h.service.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Kembalikan Response sukses beserta Token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data": LoginResponse{
			Token:       token,
			IDUser:      user.ID,
			NamaLengkap: user.NamaLengkap,
			Role:        string(user.Role),
		},
	})
}

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data user"})
	}

	// Format data agar password tidak ikut terkirim!
	var response []UserResponse
	for _, u := range users {
		response = append(response, UserResponse{
			ID:          u.ID,
			NamaLengkap: u.NamaLengkap,
			Username:    u.Username,
			Role:        string(u.Role),
			StatusAktif: u.StatusAktif,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	// Karena kita menggunakan JWT Stateless murni tanpa Database Blacklist,
	// proses "pemusnahan" token sebenarnya terjadi di sisi Frontend (React).
	// API ini bertugas memberikan respons sukses agar Frontend tahu ia boleh menghapus tokennya.
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logout berhasil. Sesi telah ditutup.",
	})
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User berhasil diupdate",
		"data": UserResponse{
			ID:      user.ID,
			NamaLengkap: user.NamaLengkap,
			Role:        string(user.Role),
		},
	})
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	// Mencegah admin menghapus dirinya sendiri secara tidak sengaja
	currentUserId := int(c.Locals("user_id").(float64)) // Diambil dari JWT
	if id == currentUserId {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Anda tidak dapat menghapus akun Anda sendiri yang sedang aktif"})
	}

	err = h.service.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User berhasil dihapus",
	})
}