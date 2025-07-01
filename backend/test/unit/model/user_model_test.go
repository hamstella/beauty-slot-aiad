package model_test

import (
	"app/src/validation"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validate = validation.Validator()

func Test_ユーザーモデル_バリデーション(t *testing.T) {
	t.Run("ユーザー作成時のバリデーション", func(t *testing.T) {
		var newUser = validation.CreateUser{
			Name:     "John Doe",
			Email:    "johndoe@gmail.com",
			Password: "password1",
			Role:     "user",
		}

		t.Run("正常なユーザーデータの場合_バリデーションが成功する", func(t *testing.T) {
			// Given: 正常なユーザーデータ（上記で定義済み）
			
			// When: バリデーションを実行
			err := validate.Struct(newUser)
			
			// Then: エラーが発生しない
			assert.NoError(t, err)
		})

		t.Run("無効なメールアドレスの場合_バリデーションエラーが発生する", func(t *testing.T) {
			// Given: 無効なメールアドレスを設定
			newUser.Email = "invalidEmail"
			
			// When: バリデーションを実行
			err := validate.Struct(newUser)
			
			// Then: バリデーションエラーが発生する
			assert.Error(t, err)
		})

		t.Run("パスワードが8文字未満の場合_バリデーションエラーが発生する", func(t *testing.T) {
			// Given: 8文字未満のパスワードを設定
			newUser.Password = "passwo1" // 7文字
			
			// When: バリデーションを実行
			err := validate.Struct(newUser)
			
			// Then: バリデーションエラーが発生する
			assert.Error(t, err)
		})

		t.Run("パスワードに数字が含まれない場合_バリデーションエラーが発生する", func(t *testing.T) {
			// Given: 数字を含まないパスワードを設定
			newUser.Password = "password" // 数字なし
			
			// When: バリデーションを実行
			err := validate.Struct(newUser)
			
			// Then: バリデーションエラーが発生する
			assert.Error(t, err)
		})

		t.Run("パスワードに文字が含まれない場合_バリデーションエラーが発生する", func(t *testing.T) {
			// Given: 文字を含まないパスワード（数字のみ）を設定
			newUser.Password = "11111111" // 数字のみ
			
			// When: バリデーションを実行
			err := validate.Struct(newUser)
			
			// Then: バリデーションエラーが発生する
			assert.Error(t, err)
		})

		t.Run("不正なロールが指定された場合_バリデーションエラーが発生する", func(t *testing.T) {
			// Given: 不正なロールを設定
			newUser.Role = "invalid" // 存在しないロール
			
			// When: バリデーションを実行
			err := validate.Struct(newUser)
			
			// Then: バリデーションエラーが発生する
			assert.Error(t, err)
		})
	})

	t.Run("ユーザー更新時のバリデーション", func(t *testing.T) {
		var updateUser = validation.UpdateUser{
			Name:     "John Doe",
			Email:    "johndoe@gmail.com",
			Password: "password1",
		}

		t.Run("should correctly validate a valid user", func(t *testing.T) {
			err := validate.Struct(updateUser)
			assert.NoError(t, err)
		})

		t.Run("should throw a validation error if email is invalid", func(t *testing.T) {
			updateUser.Email = "invalidEmail"
			err := validate.Struct(updateUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password length is less than 8 characters", func(t *testing.T) {
			updateUser.Password = "passwo1"
			err := validate.Struct(updateUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain numbers", func(t *testing.T) {
			updateUser.Password = "password"
			err := validate.Struct(updateUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain letters", func(t *testing.T) {
			updateUser.Password = "11111111"
			err := validate.Struct(updateUser)
			assert.Error(t, err)
		})
	})

	t.Run("ユーザーパスワード更新時のバリデーション", func(t *testing.T) {
		var newPassword = validation.UpdatePassOrVerify{
			Password: "password1",
		}

		t.Run("should correctly validate a valid user password", func(t *testing.T) {
			err := validate.Struct(newPassword)
			assert.NoError(t, err)
		})

		t.Run("should throw a validation error if password length is less than 8 characters", func(t *testing.T) {
			newPassword.Password = "passwo1"
			err := validate.Struct(newPassword)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain numbers", func(t *testing.T) {
			newPassword.Password = "password"
			err := validate.Struct(newPassword)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain letters", func(t *testing.T) {
			newPassword.Password = "11111111"
			err := validate.Struct(newPassword)
			assert.Error(t, err)
		})
	})

	// Note: 実際のUserモデルは存在しないため、JSON変換テストはスキップ
	// バリデーション構造体のテストのみを実行
}
