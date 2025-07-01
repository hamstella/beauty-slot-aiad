package mocks

import (
	"app/src/model"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// ReservationRepositoryMock は予約リポジトリのモック実装
type ReservationRepositoryMock struct {
	mock.Mock
}

// Create は予約を作成する
func (m *ReservationRepositoryMock) Create(reservation *model.Reservation) (*model.Reservation, error) {
	args := m.Called(reservation)
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// GetByID は予約をIDで取得する
func (m *ReservationRepositoryMock) GetByID(id uuid.UUID) (*model.Reservation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// GetByCustomerID は顧客IDで予約を取得する
func (m *ReservationRepositoryMock) GetByCustomerID(customerID uuid.UUID) ([]model.Reservation, error) {
	args := m.Called(customerID)
	return args.Get(0).([]model.Reservation), args.Error(1)
}

// GetByStaffID はスタッフIDで予約を取得する
func (m *ReservationRepositoryMock) GetByStaffID(staffID uuid.UUID) ([]model.Reservation, error) {
	args := m.Called(staffID)
	return args.Get(0).([]model.Reservation), args.Error(1)
}

// GetByDateRange は日付範囲で予約を取得する
func (m *ReservationRepositoryMock) GetByDateRange(startDate, endDate time.Time) ([]model.Reservation, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]model.Reservation), args.Error(1)
}

// Update は予約を更新する
func (m *ReservationRepositoryMock) Update(reservation *model.Reservation) (*model.Reservation, error) {
	args := m.Called(reservation)
	return args.Get(0).(*model.Reservation), args.Error(1)
}

// Delete は予約を削除する
func (m *ReservationRepositoryMock) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetConflictingReservations は時間が重複する予約を取得する
func (m *ReservationRepositoryMock) GetConflictingReservations(staffID uuid.UUID, startTime, endTime time.Time) ([]model.Reservation, error) {
	args := m.Called(staffID, startTime, endTime)
	return args.Get(0).([]model.Reservation), args.Error(1)
}

// GetPaginatedReservations はページネーション付きで予約を取得する
func (m *ReservationRepositoryMock) GetPaginatedReservations(page, limit int, filters map[string]interface{}) ([]model.Reservation, int64, error) {
	args := m.Called(page, limit, filters)
	return args.Get(0).([]model.Reservation), args.Get(1).(int64), args.Error(2)
}