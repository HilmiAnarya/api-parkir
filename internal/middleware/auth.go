package middleware

import (
	"api-parkir/internal/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 1. Middleware Validasi Token
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses ditolak, token tidak ditemukan"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secretKey := config.GetEnv("JWT_SECRET", "rahasia_parkir_2026")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau sudah kadaluarsa"})
		}

		// Simpan data token (claims) ke dalam context Fiber agar bisa dibaca oleh middleware selanjutnya
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_role", claims["role"])
		c.Locals("user_id", claims["id_user"])

		return c.Next()
	}
}

// 2. Middleware Pengecekan Role
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role").(string)

		// Cek apakah role user ada di dalam daftar role yang diizinkan
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Anda tidak memiliki hak akses untuk fitur ini"})
		}

		return c.Next()
	}
}