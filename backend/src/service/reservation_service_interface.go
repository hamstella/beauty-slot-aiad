package service

import (
	"app/src/model"

	"github.com/google/uuid"
)

// ReservationServiceInterface は予約サービスのインターフェース
type ReservationServiceInterface interface {
	GetReservations(page, limit int, status, staffID, customerID, dateFrom, dateTo string) ([]model.Reservation, int64, error)
	GetReservationByID(id uuid.UUID) (*model.Reservation, error)
	CreateReservation(reservation *model.Reservation) (*model.Reservation, error)
	UpdateReservation(reservation *model.Reservation) (*model.Reservation, error)
	CancelReservation(id uuid.UUID) error
	CreateReservationFromRequest(customerID, staffID uuid.UUID, reservationDate, startTime string, menuIDs, optionIDs []uuid.UUID, notes string) (*model.Reservation, error)
	UpdateReservationFromRequest(id, customerID, staffID uuid.UUID, reservationDate, startTime string, menuIDs, optionIDs []uuid.UUID, notes string) (*model.Reservation, error)
	UpdateReservationStatus(id uuid.UUID, status string) (*model.Reservation, error)
	GetAvailability(date, durationStr, staffIDStr, menuIDsStr string) ([]map[string]interface{}, error)
}