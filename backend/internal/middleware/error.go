package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler カスタムエラーハンドラー
func ErrorHandler(c *fiber.Ctx, err error) error {
	// デフォルトは500エラー
	code := fiber.StatusInternalServerError

	// Fiberエラーの場合はステータスコードを取得
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// エラーログ出力
	log.Printf("Error: %v", err)

	// エラーレスポンス
	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
		"code":    code,
	})
}