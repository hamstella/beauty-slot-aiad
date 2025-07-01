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
// @Summary 顧客一覧取得
// @Description ページング付きで全顧客を取得します
// @Tags 顧客管理
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "ページサイズ" default(10)
// @Success 200 {object} map[string]interface{} "顧客一覧"
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
// @Summary 顧客詳細取得
// @Description IDで指定した顧客の詳細情報を取得します
// @Tags 顧客管理
// @Accept json
// @Produce json
// @Param id path string true "顧客ID"
// @Success 200 {object} model.Customer "顧客情報"
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
// @Summary 顧客新規登録
// @Description 新しい顧客を登録します
// @Tags 顧客管理
// @Accept json
// @Produce json
// @Param customer body model.Customer true "顧客データ"
// @Success 201 {object} model.Customer "登録された顧客情報"
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
// @Summary 顧客情報更新
// @Description IDで指定した顧客の情報を更新します
// @Tags 顧客管理
// @Accept json
// @Produce json
// @Param id path string true "顧客ID"
// @Param customer body model.Customer true "更新する顧客データ"
// @Success 200 {object} model.Customer "更新された顧客情報"
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
// @Summary 顧客削除
// @Description IDで指定した顧客を削除します
// @Tags 顧客管理
// @Accept json
// @Produce json
// @Param id path string true "顧客ID"
// @Success 204 "削除成功"
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