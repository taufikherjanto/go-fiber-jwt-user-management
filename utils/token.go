package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go-fiber-user-management/model"

	"github.com/golang-jwt/jwt"
)

// GenerateToken membuat token JWT baru untuk pengguna yang ditentukan.
func GenerateToken(user model.User) (string, error) {
	// Mengambil rahasia JWT dari variabel lingkungan.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET tidak diset")
	}

	// Membuat token JWT baru dengan klaim yang mencakup ID dan email pengguna.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"issued_at": time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Menetapkan kedaluwarsa token selama 24 jam
	})

	// Menandatangani token menggunakan rahasia dan mengembalikannya.
	return token.SignedString([]byte(jwtSecret))
}

// VerifyToken memeriksa apakah token yang diberikan valid dan mengembalikan klaim.
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// Mengambil rahasia JWT dari variabel lingkungan.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET tidak diset")
	}

	// Menghapus prefix "Bearer " tanpa if statement.
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Mem-parsing dan memverifikasi token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	// Mengembalikan kesalahan jika terjadi error selama parsing.
	if err != nil {
		return nil, err
	}

	// Memeriksa apakah klaim token valid dan mengembalikannya.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token tidak valid")
}
