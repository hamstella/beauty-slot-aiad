# 美容室予約管理アプリ - 開発状況レポート

## プロジェクト概要

美容室向けの予約管理システムのMVP（Minimum Viable Product）として設計された、モバイルファースト・片手操作対応の予約アプリケーション。

### 技術スタック
- **フロントエンド**: Next.js 15 + TypeScript + Tailwind CSS 4
- **バックエンド**: Go + Fiber + GORM
- **データベース**: PostgreSQL 16
- **キャッシュ**: Redis 7
- **開発環境**: Docker Compose

## 開発方針

### 🎯 Test-Driven Development (TDD) 採用宣言

**t-wada氏の知見を基盤とした厳格なTDD実践による開発**

#### TDD基本原則
- **「テストがないコードはレガシーコード」** - 全ての新機能をテストファーストで実装
- **「テストコードもプロダクションコード」** - テストの可読性・保守性を重視
- **「Red-Green-Refactor の厳格な遵守」** - 品質と設計の継続的改善

#### テスト戦略
- **Given-When-Then構造** - テストの意図を明確化
- **日本語テスト名** - 仕様書としての機能を重視
- **AAA パターン** - Arrange-Act-Assert の徹底
- **エラーケース優先** - 異常系テストの網羅的実装

## 実装状況

### ✅ 完了済み機能

#### 開発環境・品質基盤
- **TDD環境構築完了** - testify + 継続的テスト実行環境
- Docker開発環境セットアップ完了
- PostgreSQL + Redis + pgAdmin構成
- データベース設計・マイグレーション完了
- 包括的なMakefileコマンド群（TDDワークフロー対応）

#### データベース設計・マイグレーション統合完了 ✅
- **マイグレーション統合**: ルート`database/`とボイラープレート`backend/src/database/`を統合
- **テーブル構成**: 11テーブル + 関連テーブル
  - **美容室テーブル**: customers（顧客）, staff（スタッフ）, menus（メニュー）
  - **サービステーブル**: options（オプション）, labels（ラベル）, shifts（シフト）
  - **予約テーブル**: reservations（予約）, reservation_options（予約オプション）
  - **認証テーブル**: users（ユーザー）, tokens（認証トークン）
  - **監査テーブル**: audit_logs（監査ログ）, notification_logs（通知ログ）
- **マイグレーション管理**: golang-migrate/migrateツール使用、up/downマイグレーション対応
- **セキュリティ**: UUID主キー, バリデーション制約, パーティション対応
- **監査機能**: 操作ログ・通知ログの自動記録
- **認証システム**: 管理者・スタッフ・顧客の3段階認証対応

#### バックエンドAPI基盤
- Go/Fiber APIサーバー構築完了（backend/src/main.go）
- [go-fiber-boilerplate](https://github.com/indrayyana/go-fiber-boilerplate)をベースとして使用
- RESTfulエンドポイント定義完了
  - `/api/v1/customers`, `/api/v1/reservations`, `/api/v1/health`
  - Swagger API仕様書対応（/docs）
- CORS・ログ・エラーハンドリング設定済み
- 依存関係管理（go.mod）完了
- TDDテスト基盤構築済み（testify使用）

#### フロントエンド基盤
- Next.js 15 + React 19 セットアップ完了
- Tailwind CSS 4設定済み
- TypeScript構成完了
- ESLint構成済み

### 🔄 実装中・計画中機能（TDDアプローチ）

#### APIハンドラー実装（テストファースト）
- **Red-Green-Refactorサイクル** で CRUD操作を実装
- テストケース設計 → 実装 → リファクタリングの順序で進行
- controller層のCRUD操作実装（`backend/src/controller/`）
- データベース接続・モデル実装（`backend/src/model/`）
- 空き時間検索アルゴリズム実装（アルゴリズムロジックのユニットテスト優先）

#### フロントエンド画面実装
- モバイルファースト予約フロー
  1. ホーム画面（予約履歴 + FAB）
  2. サービス選択（Bottom Sheet）
  3. 日時選択（Bottom Sheet）
  4. スタイリスト選択（Bottom Sheet）
  5. 確認・予約完了
- 管理画面（デスクトップ向け）

#### 高度な機能
- 通知システム（メール/SMS/Push）
- 認証・認可
- リアルタイム更新

## ディレクトリ構成

```
/Users/sakana/Documents/GitHub/beauty-slot-aiad/
├── backend/                    # Go/Fiber APIサーバー
│   ├── src/                   # アプリケーションロジック
│   │   ├── main.go           # エントリーポイント
│   │   ├── config/           # 設定管理
│   │   ├── controller/       # HTTPコントローラー
│   │   ├── model/            # データモデル
│   │   ├── service/          # ビジネスロジック
│   │   ├── middleware/       # ミドルウェア
│   │   ├── router/           # ルーティング
│   │   ├── database/         # DB接続・マイグレーション
│   │   ├── response/         # レスポンス構造体
│   │   ├── validation/       # バリデーション
│   │   ├── utils/            # ユーティリティ
│   │   └── docs/             # Swagger仕様書
│   ├── test/                 # テストファイル
│   │   ├── unit/             # ユニットテスト
│   │   └── mocks/            # モックファイル
│   ├── go.mod                # Go依存関係
│   ├── Dockerfile            # コンテナ設定
│   └── Makefile              # 開発コマンド
├── frontend/                  # Next.js アプリケーション
│   ├── src/app/              # App Router
│   ├── package.json          # Node.js依存関係
│   └── Dockerfile.dev        # 開発用コンテナ
├── database/migrations/       # SQLマイグレーション
│   ├── 001_create_tables.sql # テーブル定義
│   ├── 002_create_indexes.sql # インデックス
│   └── 003_insert_seed_data.sql # サンプルデータ
├── docs/                     # ドキュメント
│   ├── REQUIREMENTS.md       # 要件定義書
│   ├── data-model.md         # データモデル設計
│   └── database-schema.md    # DB設計書
├── mock/                     # UIモックアップ
│   ├── v1/                   # 初期プロトタイプ
│   └── v2/                   # 改良版プロトタイプ
├── docker-compose.yml        # 開発環境定義
├── Makefile                  # 開発コマンド
└── README.md                 # プロジェクト概要
```

## 開発環境

### クイックスタート
```bash
# 初回セットアップ
make setup

# 開発環境起動
make dev

# アクセスURL
# フロントエンド: http://localhost:3000
# バックエンドAPI: http://localhost:8080
# API仕様書: http://localhost:8080/v1/docs/
```

### 主要コマンド（Makefile）
- `make help` - 全コマンド一覧
- `make dev` - 開発環境起動
- **`make test` - テスト実行（TDD基本コマンド）**
- **`make tdd` - TDDサイクル支援（テスト監視・自動実行）**
- **`make coverage` - テストカバレッジ可視化**
- `make lint` - コード品質チェック
- `make db-reset` - データベースリセット
- `make logs` - ログ確認

## 次の開発ステップ

### 優先度: High（TDDアプローチ）
1. **APIハンドラー実装完了（テストファースト）**
   - **Red**: 失敗するテストを先に書く
   - **Green**: 最小限のコードでテストを通す
   - **Refactor**: コード品質向上と重複排除
   - CRUD操作の完全実装（各エンドポイントのテスト網羅）
   - 空き時間検索ロジック（アルゴリズムのユニットテスト優先）
   - バリデーション・エラーハンドリング（異常系テスト重視）

2. **フロントエンド予約フロー実装**
   - モバイルファースト UI コンポーネント
   - Bottom Sheet実装
   - 予約4ステップフロー

### 優先度: Medium
3. **認証システム実装**
4. **通知機能実装**
5. **管理画面実装**

### 優先度: Low
6. **テストカバレッジ向上**
7. **パフォーマンス最適化**
8. **デプロイメント準備**

## 設計の特徴

### モバイルファースト設計
- **目標**: 片手操作、最大4タップで予約完了
- **UX**: Bottom Sheet活用、親指の届く範囲でのタップ操作
- **性能要件**: JSバンドル <70KB、LCP <2.5s

### セキュリティ・監査
- HTTPS/TLS 1.3、顧客情報暗号化（AES-256）
- 全CRUD操作の監査ログ記録
- パーティション対応の高性能ログ管理

### 将来拡張性
- 決済システム連携準備
- ポイント・クーポン機能準備
- AI活用シフト最適化準備

---

**最終更新**: 2025-07-01
**開発方針**: t-wada流TDD採用による品質重視開発
**実装進捗**: TDD基盤構築完了、テストファースト実装開始段階