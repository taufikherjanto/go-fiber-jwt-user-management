package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go-fiber-user-management/database"
	"go-fiber-user-management/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Run connection to database
	database.Connect()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello this is JWT Task App")
	})

	// route authentication & task
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
