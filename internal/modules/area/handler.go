package area

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Handler struct {
	service Service 
}

func NewHandler(service Service) *Handler { 
	return &Handler{service} 
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req AreaRequest
	if err := c.BodyParser(&req); err != nil { 
		return c.Status(400).JSON(fiber.Map{"error": "Format invalid"}) 
	}
	
	area, err := h.service.CreateArea(req)

	if err != nil { return c.Status(500).JSON(fiber.Map{
		"error": err.Error()}) 
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": area})
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	areas, err := h.service.GetAllAreas()
	if err != nil { return c.Status(500).JSON(fiber.Map{
		"error": err.Error()}) 
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "data": areas})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var req AreaRequest

	if err := c.BodyParser(&req); err != nil { 
		return c.Status(400).JSON(fiber.Map{"error": "Format invalid"}) 
	}
	
	area, err := h.service.UpdateArea(uint(id), req)
	if err != nil { return c.Status(404).JSON(fiber.Map{
		"error": err.Error()}) 
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "data": area})
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteArea(uint(id)); err != nil { 
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error()}) 
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true, "message": "Area berhasil dihapus"})
}