package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"lims_auth_service/internal/database"
	"lims_auth_service/internal/handler"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	db := database.InitDB(dbName, username, password)

	jwtSecret := os.Getenv("JWT_SECRET")
	handler.SetupRoutes(app, db, jwtSecret)

	err = app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
		return
	}
}
