package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"lims_auth_service/internal/dto"
	"lims_auth_service/internal/repository"
	"lims_auth_service/internal/service"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	authService := service.NewAuthService(repo)

	api := app.Group("/api/v1")
	api.Post("/register", func(c *fiber.Ctx) error {
		var body dto.RegisterRequest

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		if body.Email == "" ||
			body.Password == "" ||
			body.FirstName == "" ||
			body.LastName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
		}

		if err := authService.Register(body); err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created"})
	})

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello, World!"})
	})
}
