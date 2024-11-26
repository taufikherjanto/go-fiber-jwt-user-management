package middleware

import (
	"strings"
	"time"

	"go-fiber-user-management/database"
	"go-fiber-user-management/model"
	"go-fiber-user-management/utils"

	"github.com/gofiber/fiber/v2"
)

// JWTAuthorization middleware memeriksa keabsahan token JWT.
// JWTAuthorization middleware for JWT authentication
func JWTAuthorization(c *fiber.Ctx) error {
	// Get token from Authorization header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Missing or malformed JWT",
		})
	}

	// Memisahkan token dari kata "Bearer"
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Authorization format",
		})
	}

	// Cek jika token telah dibatalkan di dalam basis data
	var revokedToken model.RevokedToken
	if err := database.DB.Where("token = ?", tokenString).First(&revokedToken).Error; err == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Token has been revoked",
		})
	}

	// Verify token using VerifyToken function
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid or expired token",
		})
	}

	// Check if token is expired
	if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Token has expired",
		})
	}

	// Store claims in context for later use
	c.Locals("jwt", claims)

	// If valid, proceed to the next handler
	return c.Next()
}
