package router

import (
	"app/src/controller"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReservationRoutes(api fiber.Router, db *gorm.DB) {
	reservationService := service.NewReservationService(db)
	reservationController := controller.NewReservationController(reservationService)

	reservation := api.Group("/reservations")
	reservation.Get("/", reservationController.GetReservations)
	reservation.Get("/:id", reservationController.GetReservation)
	reservation.Post("/", reservationController.CreateReservation)
	reservation.Put("/:id", reservationController.UpdateReservation)
	reservation.Put("/:id/cancel", reservationController.CancelReservation)
}