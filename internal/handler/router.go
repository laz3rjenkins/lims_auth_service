package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"lims_auth_service/internal/dto"
	"lims_auth_service/internal/middleware"
	"lims_auth_service/internal/repository"
	"lims_auth_service/internal/service"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, jwtAccessSecret, jwtRefreshSecret string) {
	repo := repository.NewUserRepository(db)
	authService := service.NewAuthService(repo, jwtAccessSecret, jwtRefreshSecret)

	api := app.Group("/api/v1")
	api.Post("/register", func(c *fiber.Ctx) error {
		var body dto.RegisterRequest

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		if body.Email == "" ||
			body.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
		}

		if err := authService.Register(body); err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User successfully created!"})
	})

	api.Post("/login", func(c *fiber.Ctx) error {
		var body struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		accessToken, refreshToken, err := authService.Login(body.Email, body.Password)
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	})

	api.Post("/refresh", func(c *fiber.Ctx) error {
		return nil
	})

	api.Get("/me", middleware.JWTProtected(jwtAccessSecret), func(c *fiber.Ctx) error {
		email := c.Locals("email").(string)
		user, err := repo.GetByEmail(email)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}

		return c.JSON(fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"is_active": user.IsActive,
		})
	})

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello, World!"})
	})
}
