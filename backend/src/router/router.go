package router

import (
	"app/src/config"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	// validate := validation.Validator()

	healthCheckService := service.NewHealthCheckService(db)
	// emailService := service.NewEmailService()
	// userService := service.NewUserService(db, validate)
	// tokenService := service.NewTokenService(db, validate, userService)
	// authService := service.NewAuthService(db, validate, userService, tokenService)

	v1 := app.Group("/v1")

	HealthCheckRoutes(v1, healthCheckService)
	// AuthRoutes(v1, authService, userService, tokenService, emailService)
	// UserRoutes(v1, userService, tokenService)
	
	// Beauty salon specific routes
	CustomerRoutes(v1, db)
	ReservationRoutes(v1, db)

	if !config.IsProd {
		DocsRoutes(v1)
	}
}
