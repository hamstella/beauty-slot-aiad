# バックエンド開発ガイド (Go/Fiber)

このバックエンドは美容室予約管理アプリのAPIサーバーです。Go + Fiber + GORM を使用し、TDD（Test-Driven Development）によって開発されています。

## 🎯 TDD開発方針

このプロジェクトでは **t-wada氏の知見** を基盤とした厳格なTDD実践を行っています。

### Red-Green-Refactor サイクル
1. **Red**: 失敗するテストを先に書く
2. **Green**: 最小限のコードでテストを通す  
3. **Refactor**: コード品質向上と重複排除

### テスト優先順位
- **エラーケース優先**: 異常系テストの網羅的実装
- **Given-When-Then構造**: テストの意図を明確化
- **日本語テスト名**: 仕様書としての機能を重視

## 🛠️ 開発環境セットアップ

### 前提条件
- Go 1.21+
- Docker & Docker Compose
- Make

### セットアップ手順
```bash
# プロジェクトルートから
make setup          # 初回セットアップ
make dev            # 開発環境起動
make test           # テスト実行
```

### バックエンド専用コマンド
```bash
# TDD開発サイクル
make tdd            # テスト監視・自動実行
make test-backend   # バックエンドテストのみ実行
make coverage       # テストカバレッジ表示

# 開発・デバッグ
make logs-backend   # バックエンドログ確認
make lint-backend   # Go linter実行
make build-backend  # バックエンドビルド

# データベース
make db-shell       # PostgreSQL接続
make db-migrate     # マイグレーション実行
make db-reset       # DB完全リセット

# API テスト
make curl-health    # ヘルスチェック
make curl-menus     # メニューAPI確認
make curl-staff     # スタッフAPI確認
```

## 🏗️ アーキテクチャ

### ディレクトリ構成
```
backend/
├── cmd/server/main.go      # エントリーポイント
├── internal/
│   ├── config/            # 設定管理
│   ├── handlers/          # HTTPハンドラー
│   ├── models/           # データモデル（GORM）
│   ├── services/         # ビジネスロジック
│   └── middleware/       # ミドルウェア
├── tests/                # テストファイル
├── go.mod               # Go依存関係
└── Dockerfile           # コンテナ設定
```

### 技術スタック詳細
- **Webフレームワーク**: [Fiber v2](https://gofiber.io/) - 高速なHTTPフレームワーク
- **ORM**: [GORM](https://gorm.io/) - Go製ORM
- **データベース**: PostgreSQL 16
- **キャッシュ**: Redis 7
- **テストフレームワーク**: [testify](https://github.com/stretchr/testify)
- **Linter**: [golangci-lint](https://golangci-lint.run/)

## 🧪 テスト戦略

### テストの種類
1. **ユニットテスト**: モデル・サービスロジック
2. **統合テスト**: API エンドポイント
3. **E2Eテスト**: 予約フロー全体

### テスト命名規約
```go
// 日本語テスト名で仕様を明確化
func Test_顧客_新規登録_正常系(t *testing.T) {
    // Given: 有効な顧客データが与えられたとき
    // When: 顧客登録APIを呼び出すと
    // Then: 顧客が正常に登録される
}

func Test_顧客_新規登録_メールアドレス重複エラー(t *testing.T) {
    // 異常系テストを重視
}
```

### テスト実行
```bash
# 継続的テスト実行
make tdd

# 単発テスト実行
make test-backend

# カバレッジ確認
make coverage
```

## 📋 開発フロー

### 新機能実装（TDDアプローチ）

1. **テストファースト実装**
   ```bash
   # テスト監視開始
   make tdd
   
   # Red: 失敗するテストを書く
   vim tests/handlers_test.go
   
   # Green: 最小限の実装
   vim internal/handlers/handlers.go
   
   # Refactor: 品質向上
   ```

2. **品質チェック**
   ```bash
   make lint-backend    # コード品質確認
   make test-backend    # 全テスト実行
   make coverage       # カバレッジ確認
   ```

3. **統合確認**
   ```bash
   make dev            # 開発環境起動
   make curl-*         # API動作確認
   ```

### コーディング規約

#### Go言語規約
- `gofmt` による自動フォーマット
- `golangci-lint` ルールに準拠
- エラーハンドリングの徹底
- 適切なパッケージ構成

#### APIデザイン規約
- RESTful原則に従う
- レスポンス形式の統一
- 適切なHTTPステータスコード
- バリデーションエラーの詳細化

#### データベース規約
- UUID主キーの使用
- 適切な制約設定
- マイグレーションによるスキーマ管理
- 監査ログの記録

## 🔧 実装ガイド

### 新しいAPIエンドポイント追加

1. **テスト作成** (`tests/handlers_test.go`)
2. **ハンドラー実装** (`internal/handlers/`)
3. **ルート登録** (`cmd/server/main.go`)
4. **モデル定義** (`internal/models/`)
5. **サービス実装** (`internal/services/`)

### データベースマイグレーション
```bash
# 新しいマイグレーションファイル作成
touch database/migrations/004_add_new_feature.sql

# マイグレーション実行
make db-migrate
```

## 🐛 デバッグ・トラブルシューティング

### ログ確認
```bash
make logs-backend    # リアルタイムログ
```

### データベース確認
```bash
make db-shell       # PostgreSQL接続
make pgadmin        # pgAdmin起動（http://localhost:5050）
```

### よくある問題
- **Port 8080 already in use**: `make down` でコンテナ停止
- **Database connection failed**: `make db-reset` でDB再構築
- **Test failures**: `make clean && make setup` で環境リセット

## 📈 品質管理

### 品質指標
- **テストカバレッジ**: 80%以上維持
- **Linter警告**: ゼロ
- **ビルド成功**: 必須

### コードレビューポイント
- TDD原則の遵守
- エラーハンドリングの適切性
- セキュリティ考慮事項
- パフォーマンス影響

## 🔗 関連リソース

- **プロジェクト概要**: [../CLAUDE.md](../CLAUDE.md)
- **全体開発ガイド**: [../CONTRIBUTING.md](../CONTRIBUTING.md)
- **データベース設計**: [../docs/database-schema.md](../docs/database-schema.md)
- **API仕様書**: http://localhost:8080/v1/docs/ (開発環境起動後)
