package main

import (
	"github.com/Akmdvtm/database"
	"github.com/Akmdvtm/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.Setup(app)
	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}
}
