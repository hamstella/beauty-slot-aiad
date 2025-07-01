package controller

import (
	"app/src/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReservationController struct {
	reservationService service.ReservationServiceInterface
}

func NewReservationController(reservationService service.ReservationServiceInterface) *ReservationController {
	return &ReservationController{
		reservationService: reservationService,
	}
}

// GetReservations godoc
// @Summary 予約一覧取得
// @Description ページング・フィルター付きで全予約を取得します
// @Tags 予約管理
// @Accept json
// @Produce json
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "ページサイズ" default(10)
// @Param status query string false "ステータス絞り込み"
// @Param date query string false "日付絞り込み (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "予約一覧"
// @Router /reservations [get]
func (c *ReservationController) GetReservations(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}
	status := ctx.Query("status")
	staffID := ctx.Query("staff_id")
	customerID := ctx.Query("customer_id")
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	reservations, total, err := c.reservationService.GetReservations(page, limit, status, staffID, customerID, dateFrom, dateTo)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "予約一覧の取得に失敗しました",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)
	hasNext := int64(page) < totalPages
	hasPrev := page > 1

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"reservations": reservations,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
				"has_next":    hasNext,
				"has_prev":    hasPrev,
			},
		},
		"meta": fiber.Map{
			"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
		},
	})
}

// GetReservation godoc
// @Summary 予約詳細取得
// @Description IDで指定した予約の詳細情報を関連データと共に取得します
// @Tags 予約管理
// @Accept json
// @Produce json
// @Param id path string true "予約ID"
// @Success 200 {object} model.Reservation "予約詳細情報"
// @Router /reservations/{id} [get]
func (c *ReservationController) GetReservation(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "無効な予約IDです",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	reservation, err := c.reservationService.GetReservationByID(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "予約が見つかりません",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"reservation": reservation,
		},
		"meta": fiber.Map{
			"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
		},
	})
}

// CreateReservation godoc
// @Summary 新規予約作成
// @Description メニューとオプションを含む新しい予約を作成します
// @Tags 予約管理
// @Accept json
// @Produce json
// @Param reservation body model.Reservation true "予約データ"
// @Success 201 {object} model.Reservation "作成された予約情報"
// @Router /reservations [post]
func (c *ReservationController) CreateReservation(ctx *fiber.Ctx) error {
	var requestBody struct {
		CustomerID      uuid.UUID   `json:"customer_id" validate:"required"`
		StaffID         uuid.UUID   `json:"staff_id" validate:"required"`
		ReservationDate string      `json:"reservation_date" validate:"required"`
		StartTime       string      `json:"start_time" validate:"required"`
		MenuIDs         []uuid.UUID `json:"menu_ids" validate:"required,min=1"`
		OptionIDs       []uuid.UUID `json:"option_ids,omitempty"`
		Notes           string      `json:"notes,omitempty" validate:"max=500"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "無効なリクエストです",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	createdReservation, err := c.reservationService.CreateReservationFromRequest(requestBody.CustomerID, requestBody.StaffID, requestBody.ReservationDate, requestBody.StartTime, requestBody.MenuIDs, requestBody.OptionIDs, requestBody.Notes)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "VALIDATION_ERROR"
		errorMsg := err.Error()
		if errorMsg == "time slot is already booked" || strings.Contains(errorMsg, "already booked") || strings.Contains(errorMsg, "重複") {
			statusCode = http.StatusConflict
			errorCode = "CONFLICT"
		}
		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    errorCode,
				"message": err.Error(),
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"reservation": createdReservation,
		},
		"meta": fiber.Map{
			"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
		},
	})
}

// UpdateReservation godoc
// @Summary 予約情報更新
// @Description IDで指定した予約の情報を更新します
// @Tags 予約管理
// @Accept json
// @Produce json
// @Param id path string true "予約ID"
// @Param reservation body model.Reservation true "更新する予約データ"
// @Success 200 {object} model.Reservation "更新された予約情報"
// @Router /reservations/{id} [put]
func (c *ReservationController) UpdateReservation(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "無効な予約IDです",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	var requestBody struct {
		CustomerID      uuid.UUID   `json:"customer_id"`
		StaffID         uuid.UUID   `json:"staff_id"`
		ReservationDate string      `json:"reservation_date"`
		StartTime       string      `json:"start_time"`
		MenuIDs         []uuid.UUID `json:"menu_ids"`
		OptionIDs       []uuid.UUID `json:"option_ids,omitempty"`
		Notes           string      `json:"notes,omitempty" validate:"max=500"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "無効なリクエストです",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	updatedReservation, err := c.reservationService.UpdateReservationFromRequest(id, requestBody.CustomerID, requestBody.StaffID, requestBody.ReservationDate, requestBody.StartTime, requestBody.MenuIDs, requestBody.OptionIDs, requestBody.Notes)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "VALIDATION_ERROR"
		if err.Error() == "reservation not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		} else if err.Error() == "cannot update cancelled or completed reservations" {
			statusCode = http.StatusUnprocessableEntity
			errorCode = "BUSINESS_RULE_ERROR"
		}
		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    errorCode,
				"message": err.Error(),
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"reservation": updatedReservation,
		},
		"meta": fiber.Map{
			"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
		},
	})
}

// CancelReservation godoc
// @Summary 予約キャンセル
// @Description IDで指定した予約をキャンセルします
// @Tags 予約管理
// @Accept json
// @Produce json
// @Param id path string true "予約ID"
// @Param reason body map[string]string true "キャンセル理由"
// @Success 200 {object} model.Reservation "キャンセル済み予約情報"
// @Router /reservations/{id}/cancel [put]
func (c *ReservationController) CancelReservation(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "無効な予約IDです",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	err = c.reservationService.CancelReservation(id)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "VALIDATION_ERROR"
		if err.Error() == "reservation not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		} else if err.Error() == "cannot cancel completed reservation" || err.Error() == "reservation is already cancelled" {
			statusCode = http.StatusUnprocessableEntity
			errorCode = "BUSINESS_RULE_ERROR"
		}
		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    errorCode,
				"message": err.Error(),
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message":        "予約をキャンセルしました",
			"reservation_id": id,
		},
		"meta": fiber.Map{
			"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
		},
	})
}

// UpdateReservationStatus godoc
// @Summary 予約ステータス更新
// @Description 予約ステータスを更新します（管理者・スタッフのみ）
// @Tags 予約管理
// @Accept json
// @Produce json
// @Param id path string true "予約ID"
// @Param status body map[string]string true "ステータス更新データ"
// @Success 200 {object} model.Reservation "更新された予約情報"
// @Router /reservations/{id}/status [patch]
func (c *ReservationController) UpdateReservationStatus(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "無効な予約IDです",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	var requestBody struct {
		Status string `json:"status" validate:"required,oneof=pending confirmed in_progress completed cancelled"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "無効なリクエストです",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	updatedReservation, err := c.reservationService.UpdateReservationStatus(id, requestBody.Status)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "VALIDATION_ERROR"
		if err.Error() == "reservation not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		} else if err.Error() == "invalid status transition" {
			statusCode = http.StatusUnprocessableEntity
			errorCode = "BUSINESS_RULE_ERROR"
		}
		return ctx.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    errorCode,
				"message": err.Error(),
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"reservation": updatedReservation,
		},
		"meta": fiber.Map{
			"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
		},
	})
}

// GetAvailability godoc
// @Summary 空き時間取得
// @Description 予約可能な時間枠を取得します
// @Tags 空き時間検索
// @Accept json
// @Produce json
// @Param date query string true "日付 (YYYY-MM-DD)"
// @Param duration query int true "所要時間（分）"
// @Param staff_id query string false "スタッフID絞り込み"
// @Param menu_ids query string false "メニューID（カンマ区切り）"
// @Success 200 {object} map[string]interface{} "空き時間一覧"
// @Router /availability [get]
func (c *ReservationController) GetAvailability(ctx *fiber.Ctx) error {
	date := ctx.Query("date")
	durationStr := ctx.Query("duration")
	staffIDStr := ctx.Query("staff_id")
	menuIDsStr := ctx.Query("menu_ids")

	if date == "" || durationStr == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": "日付と必要時間は必須です",
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	availableSlots, err := c.reservationService.GetAvailability(date, durationStr, staffIDStr, menuIDsStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
			"meta": fiber.Map{
				"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"available_slots": availableSlots,
		},
		"meta": fiber.Map{
			"timestamp": ctx.Context().Time().Format("2006-01-02T15:04:05-07:00"),
		},
	})
}