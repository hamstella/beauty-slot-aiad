# 開発ガイド

このプロジェクトは個人学習目的で作成された美容室予約管理アプリです。

## 🤖 AI開発支援について

このプロジェクトでは、生成AI（Claude）を活用した開発を行っています。
以下のガイドラインに従って、効率的で品質の高い開発を進めてください。

## 🛠️ AI支援による開発フロー

1. **要件整理**
   - 実装したい機能を明確に定義
   - 既存コードとの整合性を確認

2. **実装計画**
   - 変更が必要なファイルを特定
   - テスト戦略を立案

3. **実装・検証**
   - コード品質を保つため、テストを含めて実装
   - リントとテストを必ず実行
   ```bash
   make lint
   make test
   ```

## 🛠️ 開発環境セットアップ

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

### 開発コマンド

```bash
# ヘルプ表示
make help

# 開発環境起動/停止
make dev              # 全サービス起動
make dev-bg           # バックグラウンド起動
make down             # 停止
make clean            # 完全クリーンアップ

# データベース操作
make db-migrate       # マイグレーション実行
make db-seed          # サンプルデータ投入
make db-reset         # DB完全リセット
make db-shell         # DB接続

# ログ確認
make logs             # 全ログ
make logs-backend     # バックエンドログ
make logs-frontend    # フロントエンドログ

# テスト・品質
make test             # 全テスト実行
make lint             # コード品質チェック
make lint-fix         # 自動修正

# ビルド
make build            # 本番ビルド
```

## 🧪 テスト

```bash
# 全テスト実行
make test

# フロントエンドテスト
make test-frontend

# バックエンドテスト
make test-backend

# APIテスト（curl）
make curl-health
make curl-menus
make curl-staff
```

## 📊 開発ツール

### pgAdmin
```bash
# pgAdmin起動
make pgadmin

# アクセス: http://localhost:5050
# ID: admin@example.com / PW: admin
```

### 環境変数
```bash
# 環境変数サンプル生成
make env-example
```

## 💻 コーディング規約

### Go (Backend)
- `gofmt` でフォーマット
- `golint` でリンク
- 適切なエラーハンドリング
- テストカバレッジの維持

### TypeScript/React (Frontend)
- ESLint ルールに従う
- Prettierでフォーマット
- コンポーネントの再利用性を重視
- TypeScriptの型安全性を活用

### コミットメッセージ

このプロジェクトでは [Conventional Commits](https://www.conventionalcommits.org/) に基づいたコミットメッセージ規約を採用しています。

**基本形式:**
```
<type>(<scope>): <subject>
```

**例:**
```bash
feat(api): add reservation endpoint
fix(ui): resolve mobile responsive issue
docs(readme): update setup instructions
```

**Type例:**
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント
- `style`: フォーマット
- `refactor`: リファクタリング
- `test`: テスト
- `chore`: その他

詳細な規約は **[COMMIT_CONVENTION.md](COMMIT_CONVENTION.md)** を参照してください。

## 🐛 バグ報告

バグを発見した場合は、以下の情報を含めてIssueを作成してください：

- 期待される動作
- 実際の動作
- 再現手順
- 環境情報（OS、ブラウザなど）
- スクリーンショット（UI関連の場合）

## 💡 機能要望

新機能の提案は以下を含めてください：

- 機能の目的・背景
- 具体的な仕様
- 想定される利用シーン
- 技術的な考慮事項

## 🤖 AI開発支援のベストプラクティス

### 効果的なAI活用のコツ
- 具体的で明確な指示を出す
- 既存のコードベースとの整合性を重視
- 段階的な実装とテストを心がける
- エラーハンドリングとセキュリティを常に考慮

### AI支援による品質管理
- コードレビューの観点でAIに確認を依頼
- テストケースの網羅性をAIと検討
- リファクタリングの提案を求める
- ドキュメント生成の支援を活用
