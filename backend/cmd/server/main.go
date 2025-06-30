package main

import (
	"log"
	"os"

	"beauty-salon-reservation/internal/config"
	"beauty-salon-reservation/internal/handlers"
	"beauty-salon-reservation/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 設定読み込み
	cfg := config.Load()

	// データベース接続
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Fiberアプリケーション作成
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// ミドルウェア設定
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// ハンドラー初期化
	h := handlers.New(db)

	// ルーティング設定
	setupRoutes(app, h)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App, h *handlers.Handler) {
	// ヘルスチェック
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "beauty-salon-reservation-api",
		})
	})

	// API v1
	api := app.Group("/api/v1")

	// 顧客関連
	customers := api.Group("/customers")
	customers.Get("/", h.GetCustomers)
	customers.Post("/", h.CreateCustomer)
	customers.Get("/:id", h.GetCustomer)
	customers.Put("/:id", h.UpdateCustomer)
	customers.Delete("/:id", h.DeleteCustomer)

	// スタッフ関連
	staff := api.Group("/staff")
	staff.Get("/", h.GetStaff)
	staff.Post("/", h.CreateStaff)
	staff.Get("/:id", h.GetStaffMember)
	staff.Put("/:id", h.UpdateStaff)
	staff.Delete("/:id", h.DeleteStaff)
	staff.Get("/:id/availability", h.GetStaffAvailability)

	// メニュー関連
	menus := api.Group("/menus")
	menus.Get("/", h.GetMenus)
	menus.Post("/", h.CreateMenu)
	menus.Get("/:id", h.GetMenu)
	menus.Put("/:id", h.UpdateMenu)
	menus.Delete("/:id", h.DeleteMenu)

	// オプション関連
	options := api.Group("/options")
	options.Get("/", h.GetOptions)
	options.Post("/", h.CreateOption)
	options.Get("/:id", h.GetOption)
	options.Put("/:id", h.UpdateOption)
	options.Delete("/:id", h.DeleteOption)

	// 予約関連
	reservations := api.Group("/reservations")
	reservations.Get("/", h.GetReservations)
	reservations.Post("/", h.CreateReservation)
	reservations.Get("/:id", h.GetReservation)
	reservations.Put("/:id", h.UpdateReservation)
	reservations.Delete("/:id", h.DeleteReservation)
	reservations.Post("/:id/cancel", h.CancelReservation)

	// シフト関連
	shifts := api.Group("/shifts")
	shifts.Get("/", h.GetShifts)
	shifts.Post("/", h.CreateShift)
	shifts.Get("/:id", h.GetShift)
	shifts.Put("/:id", h.UpdateShift)
	shifts.Delete("/:id", h.DeleteShift)

	// 空き時間検索
	api.Get("/availability", h.GetAvailability)

	// ラベル関連
	labels := api.Group("/labels")
	labels.Get("/", h.GetLabels)
	labels.Post("/", h.CreateLabel)
	labels.Get("/:id", h.GetLabel)
	labels.Put("/:id", h.UpdateLabel)
	labels.Delete("/:id", h.DeleteLabel)
}