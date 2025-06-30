# コミットメッセージ規約

## 基本形式

```
<type>(<scope>): <subject>

<body>

<footer>
```

## Type（必須）

| Type | 説明 | 例 |
|------|------|-----|
| `feat` | 新機能追加 | `feat(api): add reservation endpoint` |
| `fix` | バグ修正 | `fix(auth): handle invalid token error` |
| `docs` | ドキュメント変更 | `docs(readme): update setup instructions` |
| `style` | コードフォーマット、セミコロン等 | `style(frontend): fix eslint warnings` |
| `refactor` | リファクタリング | `refactor(db): optimize reservation query` |
| `perf` | パフォーマンス改善 | `perf(api): cache menu data` |
| `test` | テスト追加・修正 | `test(reservation): add validation tests` |
| `build` | ビルドシステム、依存関係 | `build(deps): upgrade next.js to 15.3.4` |
| `ci` | CI設定変更 | `ci(github): add automated testing` |
| `chore` | その他の変更 | `chore(git): add .gitignore rules` |
| `revert` | コミットの取り消し | `revert: feat(api): add reservation endpoint` |

## Scope（推奨）

プロジェクト固有のスコープ：

| Scope | 説明 |
|-------|------|
| `api` | バックエンドAPI |
| `frontend` | フロントエンド全般 |
| `ui` | UIコンポーネント |
| `db` | データベース関連 |
| `auth` | 認証・認可 |
| `reservation` | 予約機能 |
| `staff` | スタッフ管理 |
| `customer` | 顧客管理 |
| `menu` | メニュー管理 |
| `notification` | 通知機能 |
| `config` | 設定ファイル |
| `docker` | Docker設定 |
| `deploy` | デプロイ関連 |

## Subject（必須）

- **50文字以内**で簡潔に
- 動詞で始める（add, fix, update, remove etc.）
- 現在形・命令形を使用
- 最初の文字は小文字
- 末尾にピリオドなし

### 良い例 ✅
```
feat(api): add reservation creation endpoint
fix(ui): resolve mobile responsive issue
docs(api): update authentication guide
```

### 悪い例 ❌
```
Added new feature for reservations （過去形）
Fix bug. （ピリオドあり）
update documentation （動詞が小文字で始まっていない）
feat(api): This is a very long subject line that exceeds fifty characters limit （長すぎる）
```

## Body（任意）

- subjectでは説明しきれない詳細
- **72文字で改行**
- 「何を」「なぜ」変更したかを説明
- 箇条書きOK（`-`, `*`, `+` を使用）

## Footer（任意）

- Breaking Changes: `BREAKING CHANGE: <description>`
- Issue参照: `Closes #123`, `Fixes #456`, `Refs #789`
- Co-authored-by: 共同作業者の記載

## 完全な例

```
feat(reservation): add availability search functionality

Add endpoint to search available time slots for staff members.
This enables customers to see real-time availability when booking.

- Implement /api/v1/availability endpoint
- Add staff shift validation
- Include buffer time between appointments
- Add comprehensive error handling

Closes #45
```

```
fix(auth): handle expired JWT tokens properly

Previously, expired tokens caused server crashes instead of
returning proper 401 responses.

- Add token expiration validation
- Return structured error responses
- Add proper logging for auth failures

BREAKING CHANGE: Authentication error responses now return
structured JSON instead of plain text messages.

Fixes #123
```

## 特別なケース

### マージコミット
```
Merge pull request #123 from feature/reservation-api

feat(api): add reservation management endpoints
```

### リバートコミット
```
revert: feat(auth): add OAuth integration

This reverts commit 1234567890abcdef due to security concerns.
```

### 緊急修正
```
hotfix(api): fix critical reservation data loss

Emergency fix for production issue where reservation data
was being deleted instead of updated.

Fixes #CRITICAL-789
```

## ツール設定

### Commitizen設定
```json
{
  "commitizen": {
    "path": "cz-conventional-changelog"
  }
}
```

### Pre-commit Hook例
```bash
#!/bin/sh
# .git/hooks/commit-msg

commit_regex='^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .{1,50}'

if ! grep -qE "$commit_regex" "$1"; then
    echo "Invalid commit message format!"
    echo "Expected: <type>(<scope>): <subject>"
    echo "Example: feat(api): add user authentication"
    exit 1
fi
```

## コミット頻度のガイドライン

- **Small, focused commits**: 1つの論理的変更につき1コミット
- **WIP commits**: 作業途中は `wip: working on reservation feature` 等で一時保存
- **Squash before merge**: feature branchは統合前にsquashして履歴を整理

## 日本語 vs 英語

このプロジェクトでは**日本語**でのコミットメッセージを採用します：

- チーム内での理解しやすさを重視
- 日本語での詳細な説明が可能
- プロジェクトの性質に適している

### 日本語コミットメッセージの例

```
feat(api): 予約作成エンドポイントを追加

顧客が新規予約を作成できるAPIエンドポイントを実装。
リアルタイムでの空き時間チェック機能も含む。

- /api/v1/reservations POST エンドポイント追加
- スタッフシフトとの重複チェック機能
- 予約間のバッファ時間を考慮した実装
- 包括的なエラーハンドリング

Closes #45
```

```
fix(認証): 期限切れJWTトークンのハンドリングを修正

以前は期限切れトークンがサーバークラッシュを引き起こしていたが、
適切な401レスポンスを返すように修正。

- トークン有効期限チェックを追加
- 構造化されたエラーレスポンス実装
- 認証失敗時の適切なログ出力

Fixes #123
```

---

**参考資料**
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Angular Commit Guidelines](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#commit)