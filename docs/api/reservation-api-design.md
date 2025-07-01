# 予約管理API設計書

## 概要

美容室予約管理システムのRESTful API設計仕様書です。予約の CRUD 操作を中心とした機能を提供します。

## ベースURL

```
http://localhost:8080/api/v1
```

## 認証

```
Authorization: Bearer <JWT_TOKEN>
```

## 共通レスポンス形式

### 成功レスポンス
```json
{
  "success": true,
  "data": {...},
  "meta": {
    "timestamp": "2025-06-30T10:00:00+09:00"
  }
}
```

### エラーレスポンス
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "入力データに誤りがあります",
    "details": [
      {
        "field": "phone",
        "message": "電話番号の形式が正しくありません"
      }
    ]
  },
  "meta": {
    "timestamp": "2025-06-30T10:00:00+09:00"
  }
}
```

## エンドポイント仕様

### 1. 予約一覧取得

**GET** `/reservations`

顧客の予約一覧または管理者向け全予約一覧を取得します。

#### クエリパラメータ
| パラメータ | 型 | 必須 | 説明 |
|-----------|---|------|------|
| page | integer | No | ページ番号（デフォルト: 1） |
| limit | integer | No | 1ページあたりの件数（デフォルト: 20, 最大: 100） |
| status | string | No | 予約ステータスフィルタ |
| staff_id | uuid | No | スタッフIDフィルタ |
| customer_id | uuid | No | 顧客IDフィルタ |
| date_from | date | No | 開始日フィルタ（YYYY-MM-DD） |
| date_to | date | No | 終了日フィルタ（YYYY-MM-DD） |

#### レスポンス例
```json
{
  "success": true,
  "data": {
    "reservations": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "customer_id": "550e8400-e29b-41d4-a716-446655440001",
        "staff_id": "550e8400-e29b-41d4-a716-446655440002",
        "reservation_date": "2025-07-01",
        "start_time": "10:00:00",
        "end_time": "11:30:00",
        "status": "confirmed",
        "total_price": 8000,
        "notes": "初回のお客様",
        "menus": [
          {
            "id": "550e8400-e29b-41d4-a716-446655440003",
            "name": "カット",
            "price": 4000,
            "duration": 60
          }
        ],
        "options": [
          {
            "id": "550e8400-e29b-41d4-a716-446655440004",
            "name": "シャンプー",
            "price": 1000
          }
        ],
        "customer": {
          "id": "550e8400-e29b-41d4-a716-446655440001",
          "name": "田中太郎",
          "phone": "090-1234-5678"
        },
        "staff": {
          "id": "550e8400-e29b-41d4-a716-446655440002",
          "name": "佐藤花子"
        },
        "created_at": "2025-06-29T15:00:00+09:00",
        "updated_at": "2025-06-29T15:00:00+09:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

### 2. 予約詳細取得

**GET** `/reservations/{id}`

指定した予約の詳細情報を取得します。

#### パスパラメータ
| パラメータ | 型 | 必須 | 説明 |
|-----------|---|------|------|
| id | uuid | Yes | 予約ID |

#### レスポンス例
```json
{
  "success": true,
  "data": {
    "reservation": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "customer_id": "550e8400-e29b-41d4-a716-446655440001",
      "staff_id": "550e8400-e29b-41d4-a716-446655440002",
      "reservation_date": "2025-07-01",
      "start_time": "10:00:00",
      "end_time": "11:30:00",
      "status": "confirmed",
      "total_price": 8000,
      "notes": "初回のお客様",
      "menus": [...],
      "options": [...],
      "customer": {...},
      "staff": {...},
      "created_at": "2025-06-29T15:00:00+09:00",
      "updated_at": "2025-06-29T15:00:00+09:00"
    }
  }
}
```

### 3. 予約作成

**POST** `/reservations`

新しい予約を作成します。

#### リクエストボディ
```json
{
  "customer_id": "550e8400-e29b-41d4-a716-446655440001",
  "staff_id": "550e8400-e29b-41d4-a716-446655440002",
  "reservation_date": "2025-07-01",
  "start_time": "10:00:00",
  "menu_ids": [
    "550e8400-e29b-41d4-a716-446655440003"
  ],
  "option_ids": [
    "550e8400-e29b-41d4-a716-446655440004"
  ],
  "notes": "初回のお客様"
}
```

#### バリデーションルール
- `customer_id`: 必須、有効なUUID、存在する顧客ID
- `staff_id`: 必須、有効なUUID、存在するスタッフID
- `reservation_date`: 必須、YYYY-MM-DD形式、現在日より未来の日付
- `start_time`: 必須、HH:MM:SS形式
- `menu_ids`: 必須、配列、各要素は有効なUUID
- `option_ids`: オプション、配列、各要素は有効なUUID
- `notes`: オプション、最大500文字

#### レスポンス例
```json
{
  "success": true,
  "data": {
    "reservation": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "pending",
      "total_price": 8000,
      "end_time": "11:30:00",
      ...
    }
  }
}
```

### 4. 予約更新

**PUT** `/reservations/{id}`

既存の予約を更新します。

#### パスパラメータ
| パラメータ | 型 | 必須 | 説明 |
|-----------|---|------|------|
| id | uuid | Yes | 予約ID |

#### リクエストボディ
予約作成と同じ形式ですが、部分更新をサポートします。

#### ビジネスルール
- ステータスが `completed` の予約は更新不可
- 予約時間まで2時間を切った場合は更新不可
- スタッフまたは日時変更時は空き時間チェックを実行

### 5. 予約キャンセル

**DELETE** `/reservations/{id}`

指定した予約をキャンセルします。

#### パスパラメータ
| パラメータ | 型 | 必須 | 説明 |
|-----------|---|------|------|
| id | uuid | Yes | 予約ID |

#### ビジネスルール
- ステータスが `completed` の予約はキャンセル不可
- 予約時間まで2時間を切った場合はキャンセル不可

#### レスポンス例
```json
{
  "success": true,
  "data": {
    "message": "予約をキャンセルしました",
    "reservation_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### 6. 予約ステータス更新

**PATCH** `/reservations/{id}/status`

予約のステータスのみを更新します（管理者向け）。

#### パスパラメータ
| パラメータ | 型 | 必須 | 説明 |
|-----------|---|------|------|
| id | uuid | Yes | 予約ID |

#### リクエストボディ
```json
{
  "status": "in_progress"
}
```

#### 有効なステータス遷移
- `pending` → `confirmed`, `cancelled`
- `confirmed` → `in_progress`, `cancelled`
- `in_progress` → `completed`
- `completed` → （変更不可）
- `cancelled` → （変更不可）

### 7. 空き時間検索

**GET** `/availability`

指定した条件での空き時間を検索します。

#### クエリパラメータ
| パラメータ | 型 | 必須 | 説明 |
|-----------|---|------|------|
| date | date | Yes | 検索対象日（YYYY-MM-DD） |
| duration | integer | Yes | 必要時間（分） |
| staff_id | uuid | No | 特定スタッフの空き時間のみ検索 |
| menu_ids | string | No | メニューID（カンマ区切り） |

#### レスポンス例
```json
{
  "success": true,
  "data": {
    "available_slots": [
      {
        "staff_id": "550e8400-e29b-41d4-a716-446655440002",
        "staff_name": "佐藤花子",
        "available_times": [
          {
            "start_time": "10:00:00",
            "end_time": "11:30:00"
          },
          {
            "start_time": "14:00:00",
            "end_time": "15:30:00"
          }
        ]
      }
    ]
  }
}
```

## エラーコード一覧

| コード | 説明 |
|--------|------|
| VALIDATION_ERROR | バリデーションエラー |
| NOT_FOUND | リソースが見つからない |
| CONFLICT | 重複エラー（予約時間重複等） |
| BUSINESS_RULE_ERROR | ビジネスルール違反 |
| UNAUTHORIZED | 認証エラー |
| FORBIDDEN | 認可エラー |
| INTERNAL_ERROR | システム内部エラー |

## HTTPステータスコード

| ステータス | 説明 |
|-----------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 400 | バリデーションエラー |
| 401 | 認証エラー |
| 403 | 認可エラー |
| 404 | リソースが見つからない |
| 409 | 重複エラー |
| 422 | ビジネスルール違反 |
| 500 | サーバー内部エラー |

## レスポンス時間目標

| エンドポイント | 目標時間 |
|---------------|----------|
| 予約一覧取得 | 200ms以下 |
| 予約詳細取得 | 100ms以下 |
| 予約作成 | 300ms以下 |
| 予約更新 | 300ms以下 |
| 空き時間検索 | 500ms以下 |

## セキュリティ要件

- 全エンドポイントでJWT認証を必須とする
- 顧客は自分の予約のみ参照・操作可能
- 管理者・スタッフは全予約の参照・操作が可能
- SQLインジェクション対策を実装
- レート制限を適用（1分間に60リクエスト）