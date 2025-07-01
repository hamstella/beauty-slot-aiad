package controller

import (
	"app/src/model"
	"app/src/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReservationController struct {
	reservationService *service.ReservationService
}

func NewReservationController(reservationService *service.ReservationService) *ReservationController {
	return &ReservationController{
		reservationService: reservationService,
	}
}

// GetReservations godoc
// @Summary Get all reservations
// @Description Get all reservations with pagination and filters
// @Tags reservations
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Param status query string false "Filter by status"
// @Param date query string false "Filter by date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /reservations [get]
func (c *ReservationController) GetReservations(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	status := ctx.Query("status")
	date := ctx.Query("date")

	reservations, total, err := c.reservationService.GetReservations(page, limit, status, date)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get reservations",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": reservations,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetReservation godoc
// @Summary Get reservation by ID
// @Description Get reservation by ID with related data
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID"
// @Success 200 {object} model.Reservation
// @Router /reservations/{id} [get]
func (c *ReservationController) GetReservation(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reservation ID",
		})
	}

	reservation, err := c.reservationService.GetReservationByID(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Reservation not found",
		})
	}

	return ctx.JSON(reservation)
}

// CreateReservation godoc
// @Summary Create new reservation
// @Description Create new reservation with menus and options
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body model.Reservation true "Reservation data"
// @Success 201 {object} model.Reservation
// @Router /reservations [post]
func (c *ReservationController) CreateReservation(ctx *fiber.Ctx) error {
	var reservation model.Reservation
	if err := ctx.BodyParser(&reservation); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createdReservation, err := c.reservationService.CreateReservation(&reservation)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(createdReservation)
}

// UpdateReservation godoc
// @Summary Update reservation
// @Description Update reservation by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID"
// @Param reservation body model.Reservation true "Reservation data"
// @Success 200 {object} model.Reservation
// @Router /reservations/{id} [put]
func (c *ReservationController) UpdateReservation(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reservation ID",
		})
	}

	var reservation model.Reservation
	if err := ctx.BodyParser(&reservation); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	reservation.ID = id
	updatedReservation, err := c.reservationService.UpdateReservation(&reservation)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(updatedReservation)
}

// CancelReservation godoc
// @Summary Cancel reservation
// @Description Cancel reservation by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID"
// @Param reason body map[string]string true "Cancellation reason"
// @Success 200 {object} model.Reservation
// @Router /reservations/{id}/cancel [put]
func (c *ReservationController) CancelReservation(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reservation ID",
		})
	}

	var requestBody map[string]string
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	reason := requestBody["reason"]
	cancelledReservation, err := c.reservationService.CancelReservation(id, reason)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(cancelledReservation)
}