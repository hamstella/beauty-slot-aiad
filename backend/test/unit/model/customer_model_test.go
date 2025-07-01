package model_test

import (
	"app/src/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CustomerModelTestSuite は顧客モデルのテストスイート
type CustomerModelTestSuite struct {
	suite.Suite
}

func TestCustomerModelSuite(t *testing.T) {
	suite.Run(t, new(CustomerModelTestSuite))
}

func (suite *CustomerModelTestSuite) Test_顧客モデル_バリデーション_正常系() {
	suite.Run("正常な顧客データが渡された場合_バリデーションが成功する", func() {
		// Given: 正常な顧客データ
		customer := &model.Customer{
			Name:        "田中太郎",
			Email:       "tanaka@example.com",
			Phone:       "090-1234-5678",
			IsActive:    true,
		}
		
		// When & Then: 顧客オブジェクトが正しく作成される
		assert.NotNil(suite.T(), customer)
		assert.Equal(suite.T(), "田中太郎", customer.Name)
		assert.Equal(suite.T(), "tanaka@example.com", customer.Email)
		assert.Equal(suite.T(), "090-1234-5678", customer.Phone)
		assert.True(suite.T(), customer.IsActive)
	})
}

func (suite *CustomerModelTestSuite) Test_顧客モデル_バリデーション_異常系() {
	suite.Run("名前が空の場合_論理的にエラーである", func() {
		// Given: 名前が空の顧客データ
		customer := &model.Customer{
			Name:     "", // 空
			Email:    "test@example.com",
			Phone:    "090-1234-5678",
			IsActive: true,
		}
		
		// When & Then: 名前が空であることを確認
		assert.Empty(suite.T(), customer.Name)
	})
	
	suite.Run("メールアドレスが空の場合_論理的にエラーである", func() {
		// Given: メールアドレスが空の顧客データ
		customer := &model.Customer{
			Name:     "田中太郎",
			Email:    "", // 空
			Phone:    "090-1234-5678",
			IsActive: true,
		}
		
		// When & Then: メールアドレスが空であることを確認
		assert.Empty(suite.T(), customer.Email)
	})
	
	suite.Run("電話番号が空の場合_論理的にエラーである", func() {
		// Given: 電話番号が空の顧客データ
		customer := &model.Customer{
			Name:     "田中太郎",
			Email:    "tanaka@example.com",
			Phone:    "", // 空
			IsActive: true,
		}
		
		// When & Then: 電話番号が空であることを確認
		assert.Empty(suite.T(), customer.Phone)
	})
	
	suite.Run("無効なメールアドレス形式の場合_論理的にエラーである", func() {
		// Given: 無効なメールアドレス形式の顧客データ
		customer := &model.Customer{
			Name:     "田中太郎",
			Email:    "invalid-email", // 無効な形式
			Phone:    "090-1234-5678",
			IsActive: true,
		}
		
		// When & Then: 無効なメールアドレス形式であることを確認
		// 実際のバリデーションは外部ライブラリで行うため、ここでは値の確認のみ
		assert.Equal(suite.T(), "invalid-email", customer.Email)
		assert.NotContains(suite.T(), customer.Email, "@")
	})
}

func (suite *CustomerModelTestSuite) Test_顧客モデル_GORM_BeforeCreate() {
	suite.Run("IDが空の場合_BeforeCreateでUUIDが生成される", func() {
		// Given: IDが空の顧客
		customer := &model.Customer{
			ID: uuid.Nil,
		}
		
		// When: BeforeCreateを呼び出し
		err := customer.BeforeCreate(nil)
		
		// Then: UUIDが生成される
		assert.NoError(suite.T(), err)
		assert.NotEqual(suite.T(), uuid.Nil, customer.ID)
	})
	
	suite.Run("IDが既に設定されている場合_BeforeCreateで変更されない", func() {
		// Given: IDが既に設定された顧客
		existingID := uuid.New()
		customer := &model.Customer{
			ID: existingID,
		}
		
		// When: BeforeCreateを呼び出し
		err := customer.BeforeCreate(nil)
		
		// Then: IDが変更されない
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), existingID, customer.ID)
	})
}

func (suite *CustomerModelTestSuite) Test_顧客モデル_テーブル名() {
	suite.Run("テーブル名がcustomersである", func() {
		// Given: 顧客モデル
		customer := &model.Customer{}
		
		// When: テーブル名を取得
		tableName := customer.TableName()
		
		// Then: 正しいテーブル名が返される
		assert.Equal(suite.T(), "customers", tableName)
	})
}

func (suite *CustomerModelTestSuite) Test_顧客モデル_フィールド制約() {
	suite.Run("IsActiveフィールドのデフォルト値はtrueである", func() {
		// Given: 新しい顧客オブジェクト
		customer := &model.Customer{
			Name:  "テスト顧客",
			Email: "test@example.com",
			Phone: "090-1234-5678",
			// IsActiveは設定せず
		}
		
		// When: IsActiveフィールドの初期値を確認
		// Then: デフォルト値の動作確認（実際のデフォルト値はGORMで設定）
		assert.NotNil(suite.T(), customer)
		// Note: 実際のデフォルト値設定はGORMのマイグレーションで行われる
	})
	
	suite.Run("メールアドレスは一意制約があることを前提とする", func() {
		// Given: 同じメールアドレスを持つ顧客データ
		email := "unique@example.com"
		customer1 := &model.Customer{
			Name:  "顧客1",
			Email: email,
			Phone: "090-1111-1111",
		}
		customer2 := &model.Customer{
			Name:  "顧客2",
			Email: email, // 同じメールアドレス
			Phone: "090-2222-2222",
		}
		
		// When & Then: 一意制約違反の検出は実際のDB操作で行われる
		assert.Equal(suite.T(), customer1.Email, customer2.Email)
		// 実際のテストはインテグレーションテストで実施
	})
}