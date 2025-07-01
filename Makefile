# 美容室予約管理アプリ 開発用Makefile

.PHONY: help setup dev down clean logs db-migrate db-seed test lint build

# デフォルトターゲット
help: ## このヘルプを表示
	@echo "美容室予約管理アプリ 開発コマンド"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## 初回セットアップ
	@echo "🚀 初回セットアップを開始..."
	docker-compose down -v
	docker-compose up -d postgres redis
	@echo "⏳ データベースの起動を待機中..."
	sleep 10
	docker-compose up -d backend
	sleep 5
	$(MAKE) db-migrate
	@echo "✅ セットアップ完了"

dev: ## 開発環境起動
	@echo "🏃 開発環境を起動中..."
	docker-compose up --build

dev-bg: ## 開発環境をバックグラウンドで起動
	@echo "🏃 開発環境をバックグラウンドで起動中..."
	docker-compose up -d --build

dev-frontend: ## フロントエンドのみ起動
	@echo "⚛️ フロントエンドを起動中..."
	docker-compose up frontend

dev-backend: ## バックエンドのみ起動
	@echo "🔧 バックエンドを起動中..."
	docker-compose up backend postgres redis

down: ## 開発環境停止
	@echo "⏹️ 開発環境を停止中..."
	docker-compose down

clean: ## 全てのコンテナ・ボリュームを削除
	@echo "🧹 環境をクリーンアップ中..."
	docker-compose down -v --remove-orphans
	docker system prune -f

logs: ## 全サービスのログを表示
	docker-compose logs -f

logs-backend: ## バックエンドのログを表示
	docker-compose logs -f backend

logs-frontend: ## フロントエンドのログを表示
	docker-compose logs -f frontend

logs-db: ## データベースのログを表示
	docker-compose logs -f postgres

# データベース関連
db-migrate: ## データベースマイグレーション実行
	@echo "📊 データベースマイグレーションを実行中..."
	docker-compose exec backend sh -c "cd /app && make migrate-up"

db-migrate-down: ## データベースマイグレーションロールバック
	@echo "⬇️ データベースマイグレーションをロールバック中..."
	docker-compose exec backend sh -c "cd /app && make migrate-down"

db-seed: ## サンプルデータ投入（マイグレーションに含まれています）
	@echo "🌱 サンプルデータはマイグレーションに含まれています"
	@echo "💡 db-migrateコマンドでテーブル作成とサンプルデータ投入が実行されます"

db-reset: ## データベースをリセット
	@echo "🔄 データベースをリセット中..."
	docker-compose down
	docker volume rm beauty-slot-aiad_postgres_data || true
	docker-compose up -d postgres redis
	sleep 10
	docker-compose up -d backend
	sleep 5
	$(MAKE) db-migrate

db-shell: ## データベースシェルに接続
	docker-compose exec postgres psql -U postgres -d beauty_salon_reservation

db-backup: ## データベースバックアップ
	@echo "💾 データベースをバックアップ中..."
	mkdir -p backups
	docker-compose exec -T postgres pg_dump -U postgres beauty_salon_reservation > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql

# 開発ツール
pgadmin: ## pgAdmin起動 (http://localhost:5050)
	@echo "🔧 pgAdminを起動中..."
	docker-compose --profile tools up -d pgadmin
	@echo "📊 pgAdmin: http://localhost:5050 (admin@example.com / admin)"

stop-tools: ## 開発ツール停止
	docker-compose --profile tools down

# テスト・品質・TDD
test: ## テスト実行（基本コマンド）
	@echo "🧪 テストを実行中..."
	cd frontend && npm test
	cd backend && go test ./...

test-frontend: ## フロントエンドテスト
	cd frontend && npm test

test-backend: ## バックエンドテスト
	cd backend && go test ./...

tdd: ## TDDサイクル支援（ファイル監視・自動テスト実行）
	@echo "🔄 TDDモード開始 - ファイル変更を監視してテストを自動実行"
	@echo "📝 Red-Green-Refactorサイクルでコードを書いてください"
	cd backend && find . -name "*.go" | entr -c go test ./...

tdd-verbose: ## TDD詳細モード（テスト詳細表示）
	@echo "🔍 TDD詳細モード - テスト詳細を表示"
	cd backend && find . -name "*.go" | entr -c go test -v ./...

coverage: ## テストカバレッジ表示
	@echo "📊 テストカバレッジを計算中..."
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "📋 カバレッジレポート: backend/coverage.html"

coverage-cli: ## テストカバレッジ（CLI表示）
	@echo "📊 テストカバレッジ（CLI）:"
	cd backend && go test -cover ./...

test-watch: ## テスト監視モード（変更時に自動実行）
	@echo "👀 テスト監視モード開始"
	cd backend && find . -name "*.go" | entr -c go test ./...

lint: ## コード品質チェック
	@echo "🔍 Lintを実行中..."
	cd frontend && npm run lint
	cd backend && golangci-lint run

lint-fix: ## コード自動修正
	@echo "🔧 Lintエラーを自動修正中..."
	cd frontend && npm run lint -- --fix
	cd backend && gofmt -w .

# ビルド
build: ## 本番用ビルド
	@echo "🏗️ 本番用ビルドを実行中..."
	cd frontend && npm run build
	cd backend && go build -o bin/server cmd/server/main.go

build-frontend: ## フロントエンドビルド
	cd frontend && npm run build

build-backend: ## バックエンドビルド
	cd backend && go build -o bin/server cmd/server/main.go

# API関連
api-docs: ## API仕様書を生成・表示
	@echo "📚 API仕様書を表示中..."
	@echo "Swagger UI: http://localhost:8080/swagger/"

curl-health: ## ヘルスチェック
	curl -s http://localhost:8080/health | jq

curl-menus: ## メニュー一覧取得
	curl -s http://localhost:8080/api/v1/menus | jq

curl-staff: ## スタッフ一覧取得
	curl -s http://localhost:8080/api/v1/staff | jq

# 情報表示
status: ## サービス状況確認
	@echo "📊 サービス状況:"
	docker-compose ps

urls: ## 開発用URL一覧
	@echo "🌐 開発用URL:"
	@echo "Frontend:  http://localhost:3000"
	@echo "Backend:   http://localhost:8080"
	@echo "API Docs:  http://localhost:8080/swagger/"
	@echo "pgAdmin:   http://localhost:5050 (要起動: make pgadmin)"

env-example: ## 環境変数サンプル生成
	@echo "📝 環境変数サンプルを生成中..."
	@cat > .env.example << 'EOF'
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=beauty_salon_reservation

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# API
ALLOWED_ORIGINS=http://localhost:3000
JWT_SECRET=your-secret-key-change-in-production

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
EOF
	@echo "✅ .env.example を作成しました"