package router

import (
	"app/src/controller"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CustomerRoutes(api fiber.Router, db *gorm.DB) {
	customerService := service.NewCustomerService(db)
	customerController := controller.NewCustomerController(customerService)

	customer := api.Group("/customers")
	customer.Get("/", customerController.GetCustomers)
	customer.Get("/:id", customerController.GetCustomer)
	customer.Post("/", customerController.CreateCustomer)
	customer.Put("/:id", customerController.UpdateCustomer)
	customer.Delete("/:id", customerController.DeleteCustomer)
}