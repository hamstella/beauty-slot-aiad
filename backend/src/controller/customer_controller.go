package controller

import (
	"app/src/model"
	"app/src/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CustomerController struct {
	customerService *service.CustomerService
}

func NewCustomerController(customerService *service.CustomerService) *CustomerController {
	return &CustomerController{
		customerService: customerService,
	}
}

// GetCustomers godoc
// @Summary Get all customers
// @Description Get all customers with pagination
// @Tags customers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /customers [get]
func (c *CustomerController) GetCustomers(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	customers, total, err := c.customerService.GetCustomers(page, limit)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get customers",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": customers,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetCustomer godoc
// @Summary Get customer by ID
// @Description Get customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} model.Customer
// @Router /customers/{id} [get]
func (c *CustomerController) GetCustomer(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	customer, err := c.customerService.GetCustomerByID(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Customer not found",
		})
	}

	return ctx.JSON(customer)
}

// CreateCustomer godoc
// @Summary Create new customer
// @Description Create new customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body model.Customer true "Customer data"
// @Success 201 {object} model.Customer
// @Router /customers [post]
func (c *CustomerController) CreateCustomer(ctx *fiber.Ctx) error {
	var customer model.Customer
	if err := ctx.BodyParser(&customer); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createdCustomer, err := c.customerService.CreateCustomer(&customer)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(createdCustomer)
}

// UpdateCustomer godoc
// @Summary Update customer
// @Description Update customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body model.Customer true "Customer data"
// @Success 200 {object} model.Customer
// @Router /customers/{id} [put]
func (c *CustomerController) UpdateCustomer(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	var customer model.Customer
	if err := ctx.BodyParser(&customer); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	customer.ID = id
	updatedCustomer, err := c.customerService.UpdateCustomer(&customer)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(updatedCustomer)
}

// DeleteCustomer godoc
// @Summary Delete customer
// @Description Delete customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 204
// @Router /customers/{id} [delete]
func (c *CustomerController) DeleteCustomer(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	err = c.customerService.DeleteCustomer(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Customer not found",
		})
	}

	return ctx.SendStatus(http.StatusNoContent)
}