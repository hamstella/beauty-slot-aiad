# 美容室予約管理アプリ

このプロジェクトは個人学習目的で作成された
美容室向けの予約管理システムのMVP（Minimum Viable Product）です。

## 🎯 概要

- **フロントエンド**: Next.js 15 + TypeScript + Tailwind CSS
- **バックエンド**: Go + Fiber + GORM
- **データベース**: PostgreSQL
- **キャッシュ**: Redis
- **開発環境**: Docker Compose

### 📱 UI/UX コンセプト
**モバイルファースト美容室予約アプリ** - スマートフォンで片手操作できる直感的な予約体験を提供

**設計思想**: 「数回のタップで予約完了」
- ホーム画面から最大4タップで予約完了
- Bottom Sheet を活用した階層的なUI
- 親指の届く範囲でのタップ操作
- 1画面1目的の分かりやすい情報設計

**主要画面フロー**:
1. **ホーム** - 過去の予約履歴 + FAB「＋新規予約」
2. **Step 1** - サービス選択（Bottom Sheet）
3. **Step 2** - 日時選択（Bottom Sheet）
4. **Step 3** - スタイリスト選択（Bottom Sheet）
5. **Step 4** - 確認・予約完了

詳細なUI仕様は [mock/v2/uiconcept.md](mock/v2/uiconcept.md) を参照

## 📋 機能

### 実装済み機能
- [x] データベース設計・マイグレーション
- [x] Go/Fiber API サーバー基盤
- [x] Next.js フロントエンド基盤
- [x] Docker 開発環境

### 開発予定機能
- [ ] 予約管理（CRUD操作）
- [ ] スタッフ管理・シフト管理
- [ ] メニュー・オプション管理
- [ ] 顧客管理
- [ ] 空き時間検索
- [ ] 通知機能
- [ ] 認証・認可

## 🚀 クイックスタート

### 前提条件
- Docker & Docker Compose
- Make (推奨)

### 初回セットアップ
```bash
# リポジトリクローン
git clone <repository-url>
cd beauty-salon-reservation

# 初回セットアップ（DB作成・マイグレーション・シードデータ投入）
make setup

# 開発環境起動
make dev
```

### アクセスURL
- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080
- **API仕様書**: http://localhost:8080/v1/docs/
- **データベース**: localhost:5432

## 🛠️ 開発コマンド

主要コマンド：
```bash
make help             # 全コマンド一覧
make setup            # 初回セットアップ
make dev              # 開発環境起動
make test             # テスト実行
make lint             # コード品質チェック
```

詳細は [CONTRIBUTING.md](CONTRIBUTING.md) を参照してください。

## 📁 プロジェクト構成

```
beauty-salon-reservation/
├── frontend/                 # Next.js アプリケーション
│   ├── src/
│   │   ├── app/             # App Router
│   │   ├── components/      # UIコンポーネント
│   │   └── lib/            # ユーティリティ
│   ├── package.json
│   └── Dockerfile.dev
├── backend/                  # Go/Fiber API サーバー
│   ├── cmd/server/          # エントリーポイント
│   ├── internal/
│   │   ├── handlers/        # HTTPハンドラー
│   │   ├── models/         # データモデル
│   │   ├── services/       # ビジネスロジック
│   │   ├── middleware/     # ミドルウェア
│   │   └── config/         # 設定
│   ├── go.mod
│   └── Dockerfile
├── database/                 # DB関連
│   └── migrations/         # マイグレーションSQL
├── docs/                     # ドキュメント
│   ├── REQUIREMENTS.md     # 要件定義
│   └── data/               # データ設計
│       ├── data-model.md   # データモデル
│       └── database-schema.md # DB設計書
├── docker-compose.yml       # 開発環境定義
└── Makefile                 # 開発コマンド
```

## 🗄️ データベース

### ER図
詳細は [docs/data/data-model.md](docs/data/data-model.md) を参照

### 主要テーブル
- `customers` - 顧客情報
- `staff` - スタッフ情報
- `menus` - メニュー情報
- `options` - オプション情報
- `reservations` - 予約情報
- `shifts` - シフト情報
- `labels` - ラベル（スキル分類）

## 🔧 API仕様

### エンドポイント例
```bash
# ヘルスチェック
GET /health

# メニュー一覧
GET /api/v1/menus

# スタッフ一覧
GET /api/v1/staff

# 予約作成
POST /api/v1/reservations

# 空き時間検索
GET /api/v1/availability?staff_id=xxx&date=2025-07-01&duration=60
```

詳細なAPI仕様書は開発環境起動後に http://localhost:8080/v1/docs/ で確認できます。


## 🚢 デプロイ

本番環境デプロイ手順は開発完了後に追加予定。

## 📝 要件・設計ドキュメント

- [要件定義書](docs/REQUIREMENTS.md)
- [データモデル設計](docs/data/data-model.md)
- [データベース設計書](docs/data/database-schema.md)

## 🤝 開発への参加

詳細は [CONTRIBUTING.md](CONTRIBUTING.md) を参照してください。
