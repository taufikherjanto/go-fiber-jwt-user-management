package controller

import (
	"fmt"
	"go-fiber-user-management/database"
	"go-fiber-user-management/model"
	"go-fiber-user-management/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// Register menangani pendaftaran pengguna dengan membuat pengguna baru di database.
func Register(c *fiber.Ctx) error {
	var req model.UserRequestDTO // Gunakan UserRequestDTO untuk parsing body permintaan

	// Parsing body permintaan ke dalam struct UserRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request payload",
		})
	}

	// Periksa apakah email sudah ada di database menggunakan findUserByEmail
	if _, err := findUserByEmail(req.Email); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Email already exists",
		})
	}

	// Validasi manual untuk password
	if req.Password == "" || len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Password is required and must be at least 6 characters long",
		})
	}

	// Hash password
	hashedPassword := utils.GeneratePassword(req.Password)

	// Membuat objek pengguna baru dengan data dari DTO
	user := model.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Fullname:     req.Fullname,
		Address:      req.Address,
		Gender:       req.Gender,
		PhoneNumber:  req.PhoneNumber,
	}

	// Simpan pengguna baru ke database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create user"})
	}

	// Respons sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User successfully created",
		"data":    user,
	})
}

// Login menangani login pengguna dan menghasilkan token JWT untuk pengguna yang terautentikasi.
func Login(c *fiber.Ctx) error {
	var req model.AuthenticationRequest // Gunakan model.AuthenticationRequest

	// Parsing body permintaan yang masuk ke dalam struktur authenticationRequest.
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request payload"})
	}

	// Mengambil pengguna berdasarkan email.
	user, err := findUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": err.Error()})
	}

	// Memvalidasi password yang diberikan oleh pengguna.
	if !utils.ComparePassword(user.PasswordHash, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Password salah",
		})
	}

	// Menghasilkan token JWT untuk pengguna yang terautentikasi.
	token, err := utils.GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menghasilkan token",
		})
	}

	// Mengirimkan token yang dihasilkan setelah login berhasil.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"token":   token,
	})
}

// GetUserInfo mengambil informasi pengguna berdasarkan klaim JWT.
func GetUserInfo(c *fiber.Ctx) error {
	// Mengambil klaim JWT dari konteks.
	claims := c.Locals("jwt").(jwt.MapClaims)

	// Mengambil ID pengguna dan email dari klaim
	userID, okID := claims["user_id"].(float64)
	email, okEmail := claims["email"].(string)

	if !okID || !okEmail {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Klaim token tidak valid",
		})
	}

	// Mengambil detail pengguna dari database menggunakan ID dan email.
	var user model.User
	if err := database.DB.Where("id = ? AND email = ?", userID, email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Pengguna tidak ditemukan",
		})
	}

	// Membuat instance UserResponseDTO untuk respons tanpa password.
	userResponse := model.UserResponseDTO{
		ID:          user.ID,
		Email:       user.Email,
		Fullname:    user.Fullname,
		Address:     user.Address,
		Gender:      user.Gender,
		PhoneNumber: user.PhoneNumber,
	}

	// Mengirimkan detail pengguna.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"user":    userResponse,
	})
}

// findUserByEmail mencari pengguna berdasarkan alamat email mereka.
func findUserByEmail(email string) (model.User, error) {
	var user model.User
	response := database.DB.Where("email = ?", email).First(&user)
	if response.Error != nil {
		return user, fmt.Errorf("pengguna tidak ditemukan")
	}
	return user, nil
}

// Penanganan ketika user logout
func Logout(c *fiber.Ctx) error {
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

	// Simpan token yang dibatalkan ke dalam basis data
	revokedToken := model.RevokedToken{Token: tokenString}
	if err := database.DB.Create(&revokedToken).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to revoke token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Successfully logged out",
	})
}
