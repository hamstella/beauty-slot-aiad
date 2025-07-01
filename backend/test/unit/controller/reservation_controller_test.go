package controller_test

import (
	"app/src/controller"
	"app/src/model"
	"app/test/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ReservationControllerTestSuite は予約コントローラーのテストスイート
type ReservationControllerTestSuite struct {
	suite.Suite
	app                    *fiber.App
	controller             *controller.ReservationController
	mockReservationService *mocks.ReservationServiceMock
}

func TestReservationControllerSuite(t *testing.T) {
	suite.Run(t, new(ReservationControllerTestSuite))
}

func (suite *ReservationControllerTestSuite) SetupTest() {
	// Fiberアプリとモックサービスの初期化
	suite.app = fiber.New()
	suite.mockReservationService = new(mocks.ReservationServiceMock)
	suite.controller = controller.NewReservationController(suite.mockReservationService)
	
	// ルートの設定
	suite.app.Get("/reservations", suite.controller.GetReservations)
	suite.app.Get("/reservations/:id", suite.controller.GetReservation)
	suite.app.Post("/reservations", suite.controller.CreateReservation)
	suite.app.Put("/reservations/:id", suite.controller.UpdateReservation)
	suite.app.Put("/reservations/:id/cancel", suite.controller.CancelReservation)
	suite.app.Patch("/reservations/:id/status", suite.controller.UpdateReservationStatus)
	suite.app.Get("/availability", suite.controller.GetAvailability)
}

func (suite *ReservationControllerTestSuite) TearDownTest() {
	// モックの検証
	if suite.mockReservationService != nil {
		suite.mockReservationService.AssertExpectations(suite.T())
	}
}

// エラーケース優先実装（TDDガイドライン）
func (suite *ReservationControllerTestSuite) Test_予約作成API_エラーケース() {
	suite.Run("無効なJSONが送信された場合_400_バリデーションエラーが返される", func() {
		// Given: 無効なJSONリクエスト
		invalidJSON := []byte(`{"invalid": json"}`)
		
		// When: 予約作成APIを呼び出し
		req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := suite.app.Test(req)
		
		// Then: 400エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
		
		// レスポンスボディの検証
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "VALIDATION_ERROR", response["error"].(map[string]interface{})["code"])
	})
	
	suite.Run("必須フィールドが欠けたリクエストの場合_400_バリデーションエラーが返される", func() {
		// Given: 必須フィールドが欠けたリクエスト
		invalidRequest := map[string]interface{}{
			"customer_id": "",  // 空
			"staff_id":    "",  // 空
			// reservation_date, start_time, menu_ids が欠如
		}
		
		reqBody, _ := json.Marshal(invalidRequest)
		
		// When: 予約作成APIを呼び出し
		req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := suite.app.Test(req)
		
		// Then: 400エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	})
	
	suite.Run("存在しない顧客IDが指定された場合_400_顧客が見つからないエラーが返される", func() {
		// Given: 存在しない顧客IDを含むリクエスト
		nonExistentCustomerID := uuid.New()
		validRequest := map[string]interface{}{
			"customer_id":      nonExistentCustomerID.String(),
			"staff_id":         uuid.New().String(),
			"reservation_date": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			"start_time":       "10:00:00",
			"menu_ids":         []string{uuid.New().String()},
			"notes":            "テスト予約",
		}
		
		// モックサービスの設定：顧客が見つからないエラーを返す
		suite.mockReservationService.On("CreateReservationFromRequest", 
			mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uuid.UUID"), mock.AnythingOfType("[]uuid.UUID"),
			mock.AnythingOfType("string")).
			Return(nil, fmt.Errorf("customer not found"))
		
		reqBody, _ := json.Marshal(validRequest)
		
		// When: 予約作成APIを呼び出し
		req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := suite.app.Test(req)
		
		// Then: 400エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Contains(suite.T(), response["error"].(map[string]interface{})["message"], "customer not found")
	})
	
	suite.Run("時間が重複する予約がある場合_409_競合エラーが返される", func() {
		// Given: 時間が重複するリクエスト
		conflictRequest := map[string]interface{}{
			"customer_id":      uuid.New().String(),
			"staff_id":         uuid.New().String(),
			"reservation_date": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			"start_time":       "10:00:00",
			"menu_ids":         []string{uuid.New().String()},
			"notes":            "重複テスト予約",
		}
		
		// モックサービスの設定：時間重複エラーを返す
		suite.mockReservationService.On("CreateReservationFromRequest", 
			mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uuid.UUID"), mock.AnythingOfType("[]uuid.UUID"),
			mock.AnythingOfType("string")).
			Return(nil, fmt.Errorf("time slot is already booked"))
		
		reqBody, _ := json.Marshal(conflictRequest)
		
		// When: 予約作成APIを呼び出し
		req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := suite.app.Test(req)
		
		// Then: 409エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusConflict, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "CONFLICT", response["error"].(map[string]interface{})["code"])
	})
}

func (suite *ReservationControllerTestSuite) Test_予約取得API_エラーケース() {
	suite.Run("無効なUUID形式のIDが指定された場合_400_バリデーションエラーが返される", func() {
		// Given: 無効なUUID形式のID
		invalidID := "invalid-uuid"
		
		// When: 予約取得APIを呼び出し
		req, _ := http.NewRequest("GET", "/reservations/"+invalidID, nil)
		resp, err := suite.app.Test(req)
		
		// Then: 400エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "VALIDATION_ERROR", response["error"].(map[string]interface{})["code"])
		assert.Contains(suite.T(), response["error"].(map[string]interface{})["message"], "無効な予約ID")
	})
	
	suite.Run("存在しない予約IDが指定された場合_404_予約が見つからないエラーが返される", func() {
		// Given: 存在しない予約ID
		nonExistentID := uuid.New()
		
		// モックサービスの設定：予約が見つからないエラーを返す
		suite.mockReservationService.On("GetReservationByID", nonExistentID).
			Return(nil, fmt.Errorf("reservation not found"))
		
		// When: 予約取得APIを呼び出し
		req, _ := http.NewRequest("GET", "/reservations/"+nonExistentID.String(), nil)
		resp, err := suite.app.Test(req)
		
		// Then: 404エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "NOT_FOUND", response["error"].(map[string]interface{})["code"])
	})
}

func (suite *ReservationControllerTestSuite) Test_予約更新API_エラーケース() {
	suite.Run("存在しない予約を更新しようとした場合_404_予約が見つからないエラーが返される", func() {
		// Given: 存在しない予約IDと更新データ
		nonExistentID := uuid.New()
		updateRequest := map[string]interface{}{
			"customer_id":      uuid.New().String(),
			"staff_id":         uuid.New().String(),
			"reservation_date": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			"start_time":       "11:00:00",
			"menu_ids":         []string{uuid.New().String()},
		}
		
		// モックサービスの設定：予約が見つからないエラーを返す
		suite.mockReservationService.On("UpdateReservationFromRequest",
			nonExistentID, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uuid.UUID"), mock.AnythingOfType("[]uuid.UUID"),
			mock.AnythingOfType("string")).
			Return(nil, fmt.Errorf("reservation not found"))
		
		reqBody, _ := json.Marshal(updateRequest)
		
		// When: 予約更新APIを呼び出し
		req, _ := http.NewRequest("PUT", "/reservations/"+nonExistentID.String(), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := suite.app.Test(req)
		
		// Then: 404エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "NOT_FOUND", response["error"].(map[string]interface{})["code"])
	})
	
	suite.Run("キャンセル済みの予約を更新しようとした場合_422_ビジネスルールエラーが返される", func() {
		// Given: キャンセル済み予約の更新リクエスト
		cancelledReservationID := uuid.New()
		updateRequest := map[string]interface{}{
			"notes": "更新テスト",
		}
		
		// モックサービスの設定：キャンセル済み予約の更新エラーを返す
		suite.mockReservationService.On("UpdateReservationFromRequest",
			cancelledReservationID, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uuid.UUID"), mock.AnythingOfType("[]uuid.UUID"),
			mock.AnythingOfType("string")).
			Return(nil, fmt.Errorf("cannot update cancelled or completed reservations"))
		
		reqBody, _ := json.Marshal(updateRequest)
		
		// When: 予約更新APIを呼び出し
		req, _ := http.NewRequest("PUT", "/reservations/"+cancelledReservationID.String(), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := suite.app.Test(req)
		
		// Then: 422エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusUnprocessableEntity, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "BUSINESS_RULE_ERROR", response["error"].(map[string]interface{})["code"])
	})
}

func (suite *ReservationControllerTestSuite) Test_予約キャンセルAPI_エラーケース() {
	suite.Run("存在しない予約をキャンセルしようとした場合_404_予約が見つからないエラーが返される", func() {
		// Given: 存在しない予約ID
		nonExistentID := uuid.New()
		
		// モックサービスの設定：予約が見つからないエラーを返す
		suite.mockReservationService.On("CancelReservation", nonExistentID).
			Return(fmt.Errorf("reservation not found"))
		
		// When: 予約キャンセルAPIを呼び出し
		req, _ := http.NewRequest("PUT", "/reservations/"+nonExistentID.String()+"/cancel", nil)
		resp, err := suite.app.Test(req)
		
		// Then: 404エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "NOT_FOUND", response["error"].(map[string]interface{})["code"])
	})
	
	suite.Run("既にキャンセル済みの予約をキャンセルしようとした場合_422_ビジネスルールエラーが返される", func() {
		// Given: 既にキャンセル済みの予約ID
		alreadyCancelledID := uuid.New()
		
		// モックサービスの設定：既にキャンセル済みエラーを返す
		suite.mockReservationService.On("CancelReservation", alreadyCancelledID).
			Return(fmt.Errorf("reservation is already cancelled"))
		
		// When: 予約キャンセルAPIを呼び出し
		req, _ := http.NewRequest("PUT", "/reservations/"+alreadyCancelledID.String()+"/cancel", nil)
		resp, err := suite.app.Test(req)
		
		// Then: 422エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusUnprocessableEntity, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "BUSINESS_RULE_ERROR", response["error"].(map[string]interface{})["code"])
	})
}

func (suite *ReservationControllerTestSuite) Test_ステータス更新API_エラーケース() {
	suite.Run("無効なステータス遷移を実行した場合_422_ビジネスルールエラーが返される", func() {
		// Given: 無効なステータス遷移リクエスト
		reservationID := uuid.New()
		statusRequest := map[string]interface{}{
			"status": "invalid_status",
		}
		
		// モックサービスの設定：無効な遷移エラーを返す
		suite.mockReservationService.On("UpdateReservationStatus", reservationID, "invalid_status").
			Return(nil, fmt.Errorf("invalid status transition"))
		
		reqBody, _ := json.Marshal(statusRequest)
		
		// When: ステータス更新APIを呼び出し
		req, _ := http.NewRequest("PATCH", "/reservations/"+reservationID.String()+"/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := suite.app.Test(req)
		
		// Then: 422エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusUnprocessableEntity, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "BUSINESS_RULE_ERROR", response["error"].(map[string]interface{})["code"])
	})
}

func (suite *ReservationControllerTestSuite) Test_空き時間検索API_エラーケース() {
	suite.Run("必須パラメータが欠けている場合_400_バリデーションエラーが返される", func() {
		// Given: 必須パラメータが欠けたリクエスト
		// dateパラメータなし、durationパラメータなし
		
		// When: 空き時間検索APIを呼び出し
		req, _ := http.NewRequest("GET", "/availability", nil)
		resp, err := suite.app.Test(req)
		
		// Then: 400エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Equal(suite.T(), "VALIDATION_ERROR", response["error"].(map[string]interface{})["code"])
		assert.Contains(suite.T(), response["error"].(map[string]interface{})["message"], "日付と必要時間は必須")
	})
	
	suite.Run("無効な日付形式が指定された場合_400_バリデーションエラーが返される", func() {
		// Given: 無効な日付形式のパラメータ
		invalidDate := "2024-13-32"
		duration := "60"
		
		// モックサービスの設定：無効な日付エラーを返す
		suite.mockReservationService.On("GetAvailability", invalidDate, duration, "", "").
			Return(nil, fmt.Errorf("無効な日付形式です"))
		
		// When: 空き時間検索APIを呼び出し
		req, _ := http.NewRequest("GET", "/availability?date="+invalidDate+"&duration="+duration, nil)
		resp, err := suite.app.Test(req)
		
		// Then: 400エラーが返される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(suite.T(), response["success"].(bool))
		assert.Contains(suite.T(), response["error"].(map[string]interface{})["message"], "無効な日付形式")
	})
}

// 正常系テスト（エラーケース後に実装）
func (suite *ReservationControllerTestSuite) Test_予約作成API_正常系() {
	suite.Run("正常な予約データが送信された場合_201_予約が作成される", func() {
		// Given: 正常な予約リクエスト
		customerID := uuid.New()
		staffID := uuid.New()
		validRequest := map[string]interface{}{
			"customer_id":      customerID.String(),
			"staff_id":         staffID.String(),
			"reservation_date": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			"start_time":       "10:00:00",
			"menu_ids":         []string{uuid.New().String()},
			"notes":            "正常な予約テスト",
		}
		
		// 予想される作成済み予約
		createdReservation := &model.Reservation{
			ID:         uuid.New(),
			CustomerID: customerID,
			StaffID:    staffID,
			Status:     model.ReservationStatusConfirmed,
		}
		
		// モックサービスの設定：正常な予約作成を返す
		suite.mockReservationService.On("CreateReservationFromRequest", 
			mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uuid.UUID"), mock.AnythingOfType("[]uuid.UUID"),
			mock.AnythingOfType("string")).
			Return(createdReservation, nil)
		
		reqBody, _ := json.Marshal(validRequest)
		
		// When: 予約作成APIを呼び出し
		req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := suite.app.Test(req)
		
		// Then: 201で予約が作成される
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
		
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.True(suite.T(), response["success"].(bool))
		assert.NotNil(suite.T(), response["data"].(map[string]interface{})["reservation"])
	})
}