package router

import (
	"go-fiber-user-management/controller"
	"go-fiber-user-management/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes menginisialisasi semua rute API.
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api") // Grup API utama

	// Rute Autentikasi
	auth := api.Group("/auth")            // Grup untuk rute terkait autentikasi
	auth.Post("/login", controller.Login) // Rute untuk login pengguna
	auth.Post("/register", controller.Register)
	//auth.Post("/forgot-password", controller.ForgotPassword)
	//auth.Post("/reset-password", controller.ResetPassword)                    // Rute untuk pendaftaran pengguna
	auth.Get("/profile", middleware.JWTAuthorization, controller.GetUserInfo) // Rute info pengguna yang dilindungi
	auth.Get("/logout", middleware.JWTAuthorization, controller.Logout)       // Rute info pengguna yang dilindungi

	// Route user CRUD management
	user := api.Group("/users")
	user.Get("/", middleware.JWTAuthorization, controller.GetUsers)                // Rute list pengguna oleh admin
	user.Get("/:id", middleware.JWTAuthorization, controller.GetDetailUser) // Rute untuk info pengguna oleh admin
	user.Post("/", middleware.JWTAuthorization, controller.CreateUser)
	user.Put("/:id", middleware.JWTAuthorization, controller.UpdateUser)    // Rute untuk edit pengguna oleh admin
	user.Delete("/:id", middleware.JWTAuthorization, controller.DeleteUser) // Rute untuk hapus pengguna oleh admin)
}
