# 美容室予約管理システム - バックエンドAPI

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

美容室向けの予約管理システムのMVP（Minimum Viable Product）のバックエンドAPI。Go + Fiber + PostgreSQL + Redis構成でモバイルファースト・片手操作対応の予約システムを実現。

## 特徴

- **予約管理**: 顧客・スタッフ・メニュー・予約の包括的管理
- **空き時間検索**: 高性能な空き時間検索アルゴリズム
- **監査機能**: 全操作の監査ログ記録
- **セキュリティ**: JWT認証・UUID主キー・暗号化対応
- **高性能**: PostgreSQL + Redis キャッシュ構成
- **TDD採用**: テストファースト開発による高品質実装

## クイックスタート

```bash
# プロジェクトルートから開発環境起動
cd /Users/sakana/Documents/GitHub/beauty-slot-aiad
make dev

# バックエンドのみ起動する場合
cd backend
go mod tidy
go run cmd/server/main.go
```

**アクセスURL**:
- API: http://localhost:8080
- API仕様書: http://localhost:8080/v1/docs/
- pgAdmin: http://localhost:5050

## 目次

- [API エンドポイント](#api-エンドポイント)
- [開発コマンド](#開発コマンド)
- [環境変数](#環境変数)
- [プロジェクト構成](#プロジェクト構成)
- [テスト](#テスト)
- [データベース](#データベース)
- [認証・認可](#認証認可)

## API エンドポイント

### 予約管理
- `GET /api/v1/reservations` - 予約一覧取得
- `POST /api/v1/reservations` - 新規予約作成
- `GET /api/v1/reservations/:id` - 予約詳細取得
- `PUT /api/v1/reservations/:id` - 予約更新
- `DELETE /api/v1/reservations/:id` - 予約削除

### 顧客管理
- `GET /api/v1/customers` - 顧客一覧取得
- `POST /api/v1/customers` - 新規顧客登録
- `GET /api/v1/customers/:id` - 顧客詳細取得
- `PUT /api/v1/customers/:id` - 顧客情報更新

### スタッフ管理
- `GET /api/v1/staff` - スタッフ一覧取得
- `POST /api/v1/staff` - 新規スタッフ登録
- `GET /api/v1/staff/:id` - スタッフ詳細取得
- `PUT /api/v1/staff/:id` - スタッフ情報更新

### メニュー管理
- `GET /api/v1/menus` - メニュー一覧取得
- `POST /api/v1/menus` - 新規メニュー登録
- `GET /api/v1/menus/:id` - メニュー詳細取得
- `PUT /api/v1/menus/:id` - メニュー更新

### 空き時間検索
- `GET /api/v1/availability` - 空き時間検索
  - クエリパラメータ: `date`, `staff_id`, `menu_id`

### シフト管理
- `GET /api/v1/shifts` - シフト一覧取得
- `POST /api/v1/shifts` - 新規シフト登録
- `PUT /api/v1/shifts/:id` - シフト更新

## 開発コマンド

### 基本操作
```bash
# サーバー起動
go run cmd/server/main.go

# ライブリロード（要Air）
air

# 依存関係更新
go mod tidy
```

### テスト（TDD）
```bash
# 全テスト実行
go test ./...

# テストカバレッジ
go test -cover ./...

# 特定テスト実行
go test -run TestReservationHandler ./internal/handlers

# TDDサイクル（プロジェクトルートから）
make test
make tdd
make coverage
```

### 品質管理
```bash
# リント実行
golangci-lint run

# API仕様書生成
swag init -g cmd/server/main.go -o docs
```

### データベース
```bash
# マイグレーション実行（プロジェクトルートから）
make db-reset

# 手動マイグレーション（要golang-migrate）
migrate -path database/migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
```

## 環境変数

美容室予約システム用の環境変数設定:

```bash
# サーバー設定
APP_ENV=dev
APP_HOST=0.0.0.0
APP_PORT=8080

# データベース設定（PostgreSQL）
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=beauty_salon
DB_PORT=5432

# キャッシュ設定（Redis）
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT認証
JWT_SECRET=beauty-salon-secret-key
JWT_ACCESS_EXP_MINUTES=30
JWT_REFRESH_EXP_DAYS=7

# 通知設定（将来実装）
SMTP_HOST=
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
EMAIL_FROM=noreply@beauty-salon.local
```

## プロジェクト構成

```
backend/
├── cmd/server/main.go          # エントリーポイント
├── internal/                   # アプリケーションロジック
│   ├── config/                # 設定管理
│   ├── handlers/              # HTTPハンドラー（コントローラー層）
│   ├── models/                # データモデル（データ層）
│   ├── services/              # ビジネスロジック（サービス層）
│   ├── middleware/            # カスタムミドルウェア
│   ├── utils/                 # ユーティリティ
│   └── validation/            # バリデーション
├── docs/                      # Swagger API仕様書
├── tests/                     # テストファイル
├── go.mod                     # Go依存関係
└── Dockerfile                 # コンテナ設定
```

### 主要コンポーネント
- **handlers**: REST APIエンドポイントの実装
- **models**: データベースモデル（GORM）
- **services**: 予約管理・空き時間検索等のビジネスロジック
- **middleware**: 認証・ログ・CORS等のミドルウェア

## テスト

### TDD（Test-Driven Development）採用

**t-wada氏の知見を基盤とした厳格なTDD実践**

```bash
# Red-Green-Refactorサイクル
1. Red: 失敗するテストを先に書く
2. Green: 最小限のコードでテストを通す
3. Refactor: コード品質向上と重複排除
```

### テスト実行
```bash
# 全テスト実行
go test ./...

# テストカバレッジ
go test -cover ./...

# 継続的テスト実行（プロジェクトルートから）
make tdd
```

### テスト構成
- **ユニットテスト**: handlers, services, models
- **統合テスト**: API エンドポイント
- **テストデータ**: テスト用データベース使用

## データベース

### スキーマ設計

**9テーブル + 関連テーブル構成**

- **customers**: 顧客情報（UUID主キー）
- **staff**: スタッフ情報・スキル管理
- **menus**: メニュー・料金・所要時間
- **options**: オプションメニュー
- **labels**: 分類ラベル（カテゴリ、タグ等）
- **shifts**: スタッフシフト・勤務時間
- **reservations**: 予約情報・ステータス管理
- **reservation_options**: 予約オプション関連
- **audit_logs**: 操作監査ログ（パーティション対応）
- **notification_logs**: 通知履歴

### セキュリティ
- **UUID主キー**: セキュリティ強化
- **バリデーション制約**: データ整合性保証
- **監査ログ**: 全CRUD操作の自動記録
- **暗号化**: 顧客個人情報（AES-256）

### マイグレーション
```bash
# データベースリセット（プロジェクトルートから）
make db-reset

# マイグレーションファイル場所
database/migrations/
├── 001_create_tables.sql
├── 002_create_indexes.sql
└── 003_insert_seed_data.sql
```

## 認証・認可

### JWT認証
```go
// 認証が必要なルートでの使用例
app.Post("/api/v1/reservations", middleware.Auth(), handlers.CreateReservation)
```

### 認証フロー
1. **顧客認証**: 電話番号・SMS認証（MVP）
2. **スタッフ認証**: メールアドレス・パスワード
3. **管理者認証**: 管理者専用ログイン

### 権限管理
- **顧客**: 自分の予約のみ閲覧・操作可能
- **スタッフ**: 担当予約の閲覧・更新可能
- **管理者**: 全データの閲覧・操作可能

---

## 開発状況

### ✅ 完了
- Docker開発環境構築
- データベース設計・マイグレーション
- API サーバー基盤（Go/Fiber）
- TDD環境構築

### 🔄 実装中
- APIハンドラー実装（テストファースト）
- 空き時間検索アルゴリズム
- CRUD操作完全実装

### 📋 予定
- 認証システム実装
- 通知機能（メール・SMS）
- フロントエンド連携

---

**参考**: [go-fiber-boilerplate](https://github.com/indrayyana/go-fiber-boilerplate)をベースに美容室予約システム向けにカスタマイズ
