# TDD開発ガイドライン

## t-wada流TDD実践指針

**「テストがないコードはレガシーコード」の徹底実践**

### 基本原則

#### 1. Red-Green-Refactorの厳格な遵守

```
🔴 Red    → まず失敗するテストを書く
🟢 Green  → 最小限のコードでテストを通す  
🔵 Refactor → 重複排除とコード改善
```

#### 2. テストコードもプロダクションコード

- テストの可読性・保守性を重視
- リファクタリング対象としてテストコードも継続改善
- テスト名は仕様書として機能する

#### 3. テストファースト思考

**新機能実装手順:**
1. 仕様を理解し、テストケースを設計
2. 失敗するテストを書く（Red）
3. 最小限の実装でテストを通す（Green）
4. コードを改善する（Refactor）

## テスト構造とネーミング規約

### テスト名の付け方

#### 基本パターン
```go
func Test_関数名_条件_期待結果(t *testing.T)
func Test_CreateReservation_正常な予約データが渡された場合_予約が作成される(t *testing.T)
func Test_CreateReservation_無効な日時が渡された場合_バリデーションエラーが返される(t *testing.T)
```

#### 日本語併記による仕様明確化
```go
// ✅ Good: 仕様が明確
func Test_空き時間検索_指定日時に予約がない場合_利用可能時間が返される(t *testing.T)

// ❌ Bad: 何をテストしているか不明
func TestAvailableSlots(t *testing.T)
```

### Given-When-Then構造

```go
func Test_予約作成_正常なデータ_成功(t *testing.T) {
    // Given (準備)
    repo := &MockReservationRepository{}
    service := NewReservationService(repo)
    validReservation := &Reservation{
        CustomerID: "customer-123",
        StaffID:    "staff-456",
        StartTime:  time.Now().Add(24 * time.Hour),
        MenuID:     "menu-789",
    }
    
    // When (実行)
    result, err := service.CreateReservation(validReservation)
    
    // Then (検証)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, validReservation.CustomerID, result.CustomerID)
}
```

## テスト戦略

### 1. ユニットテスト（最重要）

**対象レイヤー:**
- **Models** - ビジネスロジック、バリデーション
- **Services** - ビジネスルール、データ変換
- **Handlers** - HTTP リクエスト/レスポンス処理

**重点項目:**
- **正常系**: 仕様通りの動作確認
- **異常系**: エラーハンドリング（こちらを優先）
- **境界値**: 制限値での動作確認

### 2. インテグレーションテスト

**対象:**
- データベース操作
- 外部API連携
- エンドツーエンドの動作確認

### 3. エラーケースの優先実装

```go
// エラーケースを先にテスト
func Test_予約作成_過去の日時_エラー(t *testing.T)
func Test_予約作成_重複する予約_エラー(t *testing.T)
func Test_予約作成_存在しないスタッフID_エラー(t *testing.T)

// その後に正常系
func Test_予約作成_正常データ_成功(t *testing.T)
```

## モック・スタブ戦略

### インターフェース設計

```go
// リポジトリインターフェース
type ReservationRepository interface {
    Create(reservation *Reservation) (*Reservation, error)
    GetByID(id string) (*Reservation, error)
    Update(reservation *Reservation) (*Reservation, error)
    Delete(id string) error
}

// モック実装
type MockReservationRepository struct {
    CreateFunc func(reservation *Reservation) (*Reservation, error)
    GetByIDFunc func(id string) (*Reservation, error)
    // ...
}
```

### testifyの活用

```go
func Test_予約サービス_正常系(t *testing.T) {
    // Given
    mockRepo := &MockReservationRepository{}
    mockRepo.On("Create", mock.AnythingOfType("*Reservation")).
        Return(&Reservation{ID: "new-id"}, nil)
    
    service := NewReservationService(mockRepo)
    
    // When
    result, err := service.CreateReservation(&Reservation{})
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, "new-id", result.ID)
    mockRepo.AssertExpectations(t)
}
```

## ディレクトリ構造

```
backend/
├── internal/
│   ├── handlers/
│   │   ├── reservation_handler.go
│   │   └── reservation_handler_test.go    # ハンドラーテスト
│   ├── services/
│   │   ├── reservation_service.go
│   │   └── reservation_service_test.go    # サービステスト
│   └── models/
│       ├── reservation.go
│       └── reservation_test.go            # モデルテスト
├── test/
│   ├── integration/                       # インテグレーションテスト
│   │   ├── reservation_api_test.go
│   │   └── database_test.go
│   └── mocks/                            # モック実装
│       ├── mock_repository.go
│       └── mock_service.go
└── testdata/                             # テストデータ
    ├── fixtures/
    └── golden/
```

## TDDワークフロー

### 1. 日常的な開発サイクル

```bash
# 1. テスト監視モードで開発
make tdd

# 2. Red: 失敗するテストを書く
vim internal/services/reservation_service_test.go

# 3. Green: 最小限の実装
vim internal/services/reservation_service.go

# 4. Refactor: コード改善
# テストが通り続けることを確認しながらリファクタリング

# 5. カバレッジ確認
make coverage
```

### 2. テスト作成の順序

1. **エラーケース** から先に実装
2. **境界値テスト** で制限を明確化
3. **正常系テスト** で仕様を固定
4. **パフォーマンステスト** で品質保証

### 3. レビュー観点

#### テストコードレビュー
- [ ] テスト名が仕様を表現している
- [ ] Given-When-Then構造が明確
- [ ] エラーケースが網羅されている
- [ ] モックの使用が適切
- [ ] テストが独立している（実行順序に依存しない）

#### プロダクションコードレビュー
- [ ] 全ての新機能にテストがある
- [ ] テストが先に書かれている（git logで確認）
- [ ] エラーハンドリングが適切
- [ ] インターフェースが適切に定義されている

## よくある落とし穴と対策

### 1. テストが実装詳細に依存する

```go
// ❌ Bad: 実装詳細に依存
func Test_BadExample(t *testing.T) {
    service := NewReservationService()
    // 内部のプライベートメソッドを直接テスト
    result := service.validateDateTime(time.Now())
    assert.True(t, result)
}

// ✅ Good: 公開インターフェースをテスト
func Test_GoodExample(t *testing.T) {
    service := NewReservationService()
    // 公開メソッドの動作をテスト
    _, err := service.CreateReservation(&Reservation{})
    assert.Error(t, err) // バリデーションエラーを期待
}
```

### 2. テストが遅い

```go
// ❌ Bad: 実際のDBを使用
func Test_SlowTest(t *testing.T) {
    db := setupRealDatabase()
    defer cleanupDatabase(db)
    // ...
}

// ✅ Good: モックを使用
func Test_FastTest(t *testing.T) {
    mockRepo := &MockRepository{}
    service := NewService(mockRepo)
    // ...
}
```

### 3. テストが脆い

```go
// ❌ Bad: 固定時刻に依存
func Test_Fragile(t *testing.T) {
    reservation := &Reservation{
        StartTime: time.Date(2025, 7, 1, 10, 0, 0, 0, time.UTC),
    }
    // ...
}

// ✅ Good: 相対時刻を使用
func Test_Robust(t *testing.T) {
    reservation := &Reservation{
        StartTime: time.Now().Add(24 * time.Hour),
    }
    // ...
}
```

## メトリクス・品質指標

### カバレッジ目標
- **ユニットテスト**: 90%以上
- **統合テスト**: 70%以上
- **E2Eテスト**: 主要ユーザーフロー100%

### 品質指標
- **テスト実行時間**: ユニットテスト <10秒
- **テスト安定性**: 成功率 99.9%以上
- **メンテナンス性**: テストコードの行数 < プロダクションコードの1.5倍

---

**参考資料:**
- [t-wada氏のTDD資料](https://www.slideshare.net/t_wada/tdd-16471474)
- [テスト駆動開発](https://www.amazon.co.jp/dp/4274217884)
- [Go testing パッケージ](https://pkg.go.dev/testing)
- [testify ライブラリ](https://github.com/stretchr/testify)

**最終更新**: 2025-07-01  
**適用プロジェクト**: 美容室予約管理アプリ