package mocks

import (
	"app/src/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// ReservationServiceMock は予約サービスのモック実装
type ReservationServiceMock struct {
	mock.Mock
}

// GetReservations は予約一覧を取得する
func (m *ReservationServiceMock) GetReservations(page, limit int, status, staffID, customerID, dateFrom, dateTo string) ([]model.Reservation, int64, error) {
	args := m.Called(page, limit, status, staffID, customerID, dateFrom, dateTo)
	return args.Get(0).([]model.Reservation), args.Get(1).(int64), args.Error(2)
}

// GetReservationByID は予約をIDで取得する
func (m *ReservationServiceMock) GetReservationByID(id uuid.UUID) (*model.Reservation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// CreateReservation は予約を作成する
func (m *ReservationServiceMock) CreateReservation(reservation *model.Reservation) (*model.Reservation, error) {
	args := m.Called(reservation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// UpdateReservation は予約を更新する
func (m *ReservationServiceMock) UpdateReservation(reservation *model.Reservation) (*model.Reservation, error) {
	args := m.Called(reservation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// CancelReservation は予約をキャンセルする
func (m *ReservationServiceMock) CancelReservation(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// CreateReservationFromRequest はリクエストデータから予約を作成する
func (m *ReservationServiceMock) CreateReservationFromRequest(customerID, staffID uuid.UUID, reservationDate, startTime string, menuIDs, optionIDs []uuid.UUID, notes string) (*model.Reservation, error) {
	args := m.Called(customerID, staffID, reservationDate, startTime, menuIDs, optionIDs, notes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// UpdateReservationFromRequest はリクエストデータから予約を更新する
func (m *ReservationServiceMock) UpdateReservationFromRequest(id, customerID, staffID uuid.UUID, reservationDate, startTime string, menuIDs, optionIDs []uuid.UUID, notes string) (*model.Reservation, error) {
	args := m.Called(id, customerID, staffID, reservationDate, startTime, menuIDs, optionIDs, notes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// UpdateReservationStatus は予約ステータスを更新する
func (m *ReservationServiceMock) UpdateReservationStatus(id uuid.UUID, status string) (*model.Reservation, error) {
	args := m.Called(id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// GetAvailability は空き時間を検索する
func (m *ReservationServiceMock) GetAvailability(date, durationStr, staffIDStr, menuIDsStr string) ([]map[string]interface{}, error) {
	args := m.Called(date, durationStr, staffIDStr, menuIDsStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}