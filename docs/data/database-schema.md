# PostgreSQL テーブル定義書

## 1. 概要

美容室予約管理アプリ用のPostgreSQLデータベーススキーマ定義書。
Go/Fiber + GORM での利用を想定した設計となっている。

## 2. データベース設定

### 2.1 基本設定
```sql
-- データベース作成
CREATE DATABASE beauty_salon_reservation 
WITH 
    ENCODING = 'UTF8'
    LC_COLLATE = 'ja_JP.UTF-8'
    LC_CTYPE = 'ja_JP.UTF-8'
    TEMPLATE = template0;

-- タイムゾーン設定
SET timezone = 'Asia/Tokyo';

-- UUID拡張有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

### 2.2 接続プール推奨設定
```
max_connections = 100
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB
```

## 3. テーブル定義

### 3.1 customers（顧客）

```sql
CREATE TABLE customers (
    customer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT customers_name_check CHECK (LENGTH(TRIM(name)) > 0),
    CONSTRAINT customers_phone_check CHECK (phone ~ '^[0-9\-\+\(\)\s]+$'),
    CONSTRAINT customers_email_check CHECK (email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- インデックス
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_created_at ON customers(created_at);

-- 更新日時自動更新トリガー
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 3.2 staff（スタッフ）

```sql
CREATE TABLE staff (
    staff_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT staff_name_check CHECK (LENGTH(TRIM(name)) > 0)
);

-- インデックス
CREATE INDEX idx_staff_is_active ON staff(is_active);
CREATE INDEX idx_staff_name ON staff(name);

-- 更新日時自動更新トリガー
CREATE TRIGGER update_staff_updated_at BEFORE UPDATE ON staff
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 3.3 menus（メニュー）

```sql
CREATE TABLE menus (
    menu_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    duration_minutes INTEGER NOT NULL,
    price INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT menus_name_check CHECK (LENGTH(TRIM(name)) > 0),
    CONSTRAINT menus_duration_check CHECK (duration_minutes > 0),
    CONSTRAINT menus_price_check CHECK (price >= 0)
);

-- インデックス
CREATE INDEX idx_menus_is_active ON menus(is_active);
CREATE INDEX idx_menus_price ON menus(price);
CREATE INDEX idx_menus_duration ON menus(duration_minutes);

-- 全文検索インデックス
CREATE INDEX idx_menus_name_fulltext ON menus USING gin(to_tsvector('japanese', name));

-- 更新日時自動更新トリガー
CREATE TRIGGER update_menus_updated_at BEFORE UPDATE ON menus
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 3.4 options（オプション）

```sql
CREATE TABLE options (
    option_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    add_duration_minutes INTEGER NOT NULL,
    add_price INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT options_name_check CHECK (LENGTH(TRIM(name)) > 0),
    CONSTRAINT options_duration_check CHECK (add_duration_minutes >= 0),
    CONSTRAINT options_price_check CHECK (add_price >= 0)
);

-- インデックス
CREATE INDEX idx_options_is_active ON options(is_active);
CREATE INDEX idx_options_add_price ON options(add_price);

-- 全文検索インデックス
CREATE INDEX idx_options_name_fulltext ON options USING gin(to_tsvector('japanese', name));

-- 更新日時自動更新トリガー
CREATE TRIGGER update_options_updated_at BEFORE UPDATE ON options
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 3.5 labels（ラベル）

```sql
CREATE TABLE labels (
    label_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT labels_name_check CHECK (LENGTH(TRIM(name)) > 0)
);

-- インデックス
CREATE INDEX idx_labels_name ON labels(name);

-- 更新日時自動更新トリガー
CREATE TRIGGER update_labels_updated_at BEFORE UPDATE ON labels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 3.6 staff_labels（スタッフラベル関連）

```sql
CREATE TABLE staff_labels (
    staff_id UUID NOT NULL,
    label_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (staff_id, label_id),
    FOREIGN KEY (staff_id) REFERENCES staff(staff_id) ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES labels(label_id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX idx_staff_labels_staff_id ON staff_labels(staff_id);
CREATE INDEX idx_staff_labels_label_id ON staff_labels(label_id);
```

### 3.7 menu_labels（メニューラベル関連）

```sql
CREATE TABLE menu_labels (
    menu_id UUID NOT NULL,
    label_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (menu_id, label_id),
    FOREIGN KEY (menu_id) REFERENCES menus(menu_id) ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES labels(label_id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX idx_menu_labels_menu_id ON menu_labels(menu_id);
CREATE INDEX idx_menu_labels_label_id ON menu_labels(label_id);
```

### 3.8 shifts（シフト）

```sql
CREATE TABLE shifts (
    shift_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL,
    work_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (staff_id) REFERENCES staff(staff_id) ON DELETE CASCADE,
    UNIQUE (staff_id, work_date),
    CONSTRAINT shifts_time_check CHECK (end_time > start_time)
);

-- インデックス
CREATE INDEX idx_shifts_staff_date ON shifts(staff_id, work_date);
CREATE INDEX idx_shifts_work_date ON shifts(work_date);
CREATE INDEX idx_shifts_staff_id ON shifts(staff_id);

-- 更新日時自動更新トリガー
CREATE TRIGGER update_shifts_updated_at BEFORE UPDATE ON shifts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 3.9 reservations（予約）

```sql
CREATE TABLE reservations (
    reservation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    staff_id UUID NOT NULL,
    menu_id UUID NOT NULL,
    start_at TIMESTAMP WITH TIME ZONE NOT NULL,
    end_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_price INTEGER NOT NULL,
    total_duration_minutes INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE RESTRICT,
    FOREIGN KEY (staff_id) REFERENCES staff(staff_id) ON DELETE RESTRICT,
    FOREIGN KEY (menu_id) REFERENCES menus(menu_id) ON DELETE RESTRICT,
    
    CONSTRAINT reservations_time_check CHECK (end_at > start_at),
    CONSTRAINT reservations_status_check CHECK (status IN ('pending', 'confirmed', 'in_progress', 'completed', 'cancelled')),
    CONSTRAINT reservations_price_check CHECK (total_price >= 0),
    CONSTRAINT reservations_duration_check CHECK (total_duration_minutes > 0)
);

-- インデックス（重要なクエリパターンに最適化）
CREATE INDEX idx_reservations_staff_start_at ON reservations(staff_id, start_at);
CREATE INDEX idx_reservations_customer_id ON reservations(customer_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_start_at ON reservations(start_at);
CREATE INDEX idx_reservations_staff_status ON reservations(staff_id, status);

-- 重複予約防止制約（同一スタッフ・同一時間帯）
CREATE UNIQUE INDEX idx_reservations_no_overlap 
ON reservations(staff_id, start_at, end_at) 
WHERE status NOT IN ('cancelled');

-- 更新日時自動更新トリガー
CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 3.10 reservation_options（予約オプション関連）

```sql
CREATE TABLE reservation_options (
    reservation_id UUID NOT NULL,
    option_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (reservation_id, option_id),
    FOREIGN KEY (reservation_id) REFERENCES reservations(reservation_id) ON DELETE CASCADE,
    FOREIGN KEY (option_id) REFERENCES options(option_id) ON DELETE RESTRICT
);

-- インデックス
CREATE INDEX idx_reservation_options_reservation_id ON reservation_options(reservation_id);
CREATE INDEX idx_reservation_options_option_id ON reservation_options(option_id);
```

### 3.11 audit_logs（監査ログ）- パーティション対応

```sql
CREATE TABLE audit_logs (
    log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(20) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT audit_logs_action_check CHECK (action IN ('CREATE', 'UPDATE', 'DELETE')),
    CONSTRAINT audit_logs_entity_type_check CHECK (LENGTH(TRIM(entity_type)) > 0)
) PARTITION BY RANGE (created_at);

-- 月次パーティション作成（例：2025年）
CREATE TABLE audit_logs_2025_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE audit_logs_2025_02 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
-- 以下、必要に応じて作成

-- インデックス
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- JSONB インデックス（検索性能向上）
CREATE INDEX idx_audit_logs_old_values ON audit_logs USING gin(old_values);
CREATE INDEX idx_audit_logs_new_values ON audit_logs USING gin(new_values);
```

### 3.12 notification_logs（通知ログ）- パーティション対応

```sql
CREATE TABLE notification_logs (
    notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reservation_id UUID,
    notification_type VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    message_content TEXT,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (reservation_id) REFERENCES reservations(reservation_id) ON DELETE SET NULL,
    CONSTRAINT notification_logs_channel_check CHECK (channel IN ('email', 'sms', 'push')),
    CONSTRAINT notification_logs_status_check CHECK (status IN ('pending', 'sent', 'failed', 'delivered'))
) PARTITION BY RANGE (created_at);

-- 月次パーティション作成（例：2025年）
CREATE TABLE notification_logs_2025_01 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE notification_logs_2025_02 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
-- 以下、必要に応じて作成

-- インデックス
CREATE INDEX idx_notification_logs_reservation_id ON notification_logs(reservation_id);
CREATE INDEX idx_notification_logs_status ON notification_logs(status);
CREATE INDEX idx_notification_logs_sent_at ON notification_logs(sent_at);
CREATE INDEX idx_notification_logs_channel ON notification_logs(channel);
CREATE INDEX idx_notification_logs_recipient ON notification_logs(recipient);
```

## 4. ビュー定義

### 4.1 予約詳細ビュー

```sql
CREATE VIEW reservation_details AS
SELECT 
    r.reservation_id,
    r.start_at,
    r.end_at,
    r.status,
    r.total_price,
    r.total_duration_minutes,
    c.name AS customer_name,
    c.phone AS customer_phone,
    c.email AS customer_email,
    s.name AS staff_name,
    m.name AS menu_name,
    m.price AS menu_price,
    m.duration_minutes AS menu_duration,
    COALESCE(
        json_agg(
            json_build_object(
                'option_id', opt.option_id,
                'name', opt.name,
                'add_price', opt.add_price,
                'add_duration_minutes', opt.add_duration_minutes
            )
        ) FILTER (WHERE opt.option_id IS NOT NULL), 
        '[]'::json
    ) AS options,
    r.created_at,
    r.updated_at
FROM reservations r
JOIN customers c ON r.customer_id = c.customer_id
JOIN staff s ON r.staff_id = s.staff_id
JOIN menus m ON r.menu_id = m.menu_id
LEFT JOIN reservation_options ro ON r.reservation_id = ro.reservation_id
LEFT JOIN options opt ON ro.option_id = opt.option_id
GROUP BY 
    r.reservation_id, r.start_at, r.end_at, r.status, r.total_price, r.total_duration_minutes,
    c.name, c.phone, c.email, s.name, m.name, m.price, m.duration_minutes,
    r.created_at, r.updated_at;
```

### 4.2 スタッフ空き時間ビュー

```sql
CREATE VIEW staff_availability AS
SELECT 
    s.staff_id,
    s.name AS staff_name,
    sh.work_date,
    sh.start_time,
    sh.end_time,
    COALESCE(
        json_agg(
            json_build_object(
                'reservation_id', r.reservation_id,
                'start_at', r.start_at,
                'end_at', r.end_at,
                'status', r.status
            )
            ORDER BY r.start_at
        ) FILTER (WHERE r.reservation_id IS NOT NULL), 
        '[]'::json
    ) AS reservations
FROM staff s
JOIN shifts sh ON s.staff_id = sh.staff_id
LEFT JOIN reservations r ON s.staff_id = r.staff_id 
    AND DATE(r.start_at AT TIME ZONE 'Asia/Tokyo') = sh.work_date
    AND r.status NOT IN ('cancelled')
WHERE s.is_active = true
GROUP BY s.staff_id, s.name, sh.work_date, sh.start_time, sh.end_time;
```

## 5. 関数定義

### 5.1 空き時間検索関数

```sql
CREATE OR REPLACE FUNCTION find_available_slots(
    p_staff_id UUID,
    p_date DATE,
    p_duration_minutes INTEGER
)
RETURNS TABLE(
    available_start_time TIMESTAMP WITH TIME ZONE,
    available_end_time TIMESTAMP WITH TIME ZONE
) AS $$
DECLARE
    shift_start TIME;
    shift_end TIME;
    slot_start TIMESTAMP WITH TIME ZONE;
    slot_end TIMESTAMP WITH TIME ZONE;
    next_reservation_start TIMESTAMP WITH TIME ZONE;
BEGIN
    -- スタッフのシフト情報取得
    SELECT start_time, end_time INTO shift_start, shift_end
    FROM shifts 
    WHERE staff_id = p_staff_id AND work_date = p_date;
    
    IF shift_start IS NULL THEN
        RETURN; -- シフトなし
    END IF;
    
    -- 開始時刻を設定
    slot_start := (p_date + shift_start) AT TIME ZONE 'Asia/Tokyo';
    
    -- 予約済み時間を避けて空き時間を検索
    FOR next_reservation_start IN 
        SELECT start_at 
        FROM reservations 
        WHERE staff_id = p_staff_id 
            AND DATE(start_at AT TIME ZONE 'Asia/Tokyo') = p_date
            AND status NOT IN ('cancelled')
        ORDER BY start_at
    LOOP
        slot_end := slot_start + (p_duration_minutes || ' minutes')::INTERVAL;
        
        -- 空き時間が予約時間と重複しないかチェック
        IF slot_end <= next_reservation_start THEN
            available_start_time := slot_start;
            available_end_time := slot_end;
            RETURN NEXT;
        END IF;
        
        -- 次の開始時刻を予約終了時刻に設定
        SELECT end_at INTO slot_start
        FROM reservations 
        WHERE staff_id = p_staff_id 
            AND start_at = next_reservation_start;
    END LOOP;
    
    -- 最後の予約後の空き時間をチェック
    slot_end := slot_start + (p_duration_minutes || ' minutes')::INTERVAL;
    IF slot_end <= (p_date + shift_end) AT TIME ZONE 'Asia/Tokyo' THEN
        available_start_time := slot_start;
        available_end_time := slot_end;
        RETURN NEXT;
    END IF;
    
    RETURN;
END;
$$ LANGUAGE plpgsql;
```

## 6. セキュリティ設定

### 6.1 行レベルセキュリティ（RLS）

```sql
-- 顧客テーブルのRLS設定例
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

-- 顧客は自分の情報のみ参照可能
CREATE POLICY customer_select_policy ON customers
    FOR SELECT USING (customer_id = current_setting('app.current_customer_id')::UUID);

-- 予約テーブルのRLS設定例
ALTER TABLE reservations ENABLE ROW LEVEL SECURITY;

-- 顧客は自分の予約のみ参照可能
CREATE POLICY customer_reservation_policy ON reservations
    FOR SELECT USING (customer_id = current_setting('app.current_customer_id')::UUID);
```

### 6.2 ロール・権限設定

```sql
-- アプリケーション用ロール作成
CREATE ROLE app_api_user WITH LOGIN PASSWORD 'secure_password_here';
CREATE ROLE app_readonly_user WITH LOGIN PASSWORD 'readonly_password_here';

-- 基本権限付与
GRANT CONNECT ON DATABASE beauty_salon_reservation TO app_api_user;
GRANT USAGE ON SCHEMA public TO app_api_user;

-- テーブル権限付与
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_api_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_api_user;

-- 読み取り専用ユーザー
GRANT CONNECT ON DATABASE beauty_salon_reservation TO app_readonly_user;
GRANT USAGE ON SCHEMA public TO app_readonly_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO app_readonly_user;
```

## 7. パフォーマンス最適化

### 7.1 統計情報更新スケジュール

```sql
-- 統計情報の自動更新設定
ALTER TABLE reservations SET (autovacuum_analyze_scale_factor = 0.02);
ALTER TABLE audit_logs SET (autovacuum_analyze_scale_factor = 0.1);
```

### 7.2 パーティション管理

```sql
-- パーティション自動作成関数
CREATE OR REPLACE FUNCTION create_monthly_partitions(
    table_name TEXT,
    start_date DATE,
    num_months INTEGER
)
RETURNS VOID AS $$
DECLARE
    partition_date DATE;
    partition_name TEXT;
    next_partition_date DATE;
BEGIN
    FOR i IN 0..num_months-1 LOOP
        partition_date := start_date + (i || ' months')::INTERVAL;
        next_partition_date := partition_date + '1 month'::INTERVAL;
        partition_name := table_name || '_' || to_char(partition_date, 'YYYY_MM');
        
        EXECUTE format('CREATE TABLE IF NOT EXISTS %I PARTITION OF %I FOR VALUES FROM (%L) TO (%L)',
            partition_name, table_name, partition_date, next_partition_date);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- 使用例：6ヶ月分のパーティションを作成
SELECT create_monthly_partitions('audit_logs', '2025-01-01'::DATE, 6);
SELECT create_monthly_partitions('notification_logs', '2025-01-01'::DATE, 6);
```

## 8. バックアップ・メンテナンス

### 8.1 バックアップスクリプト例

```bash
#!/bin/bash
# backup_database.sh

DB_NAME="beauty_salon_reservation"
BACKUP_DIR="/var/backups/postgresql"
DATE=$(date +%Y%m%d_%H%M%S)

# フルバックアップ
pg_dump -U postgres -h localhost -d $DB_NAME > "$BACKUP_DIR/full_backup_$DATE.sql"

# 圧縮
gzip "$BACKUP_DIR/full_backup_$DATE.sql"

# 古いバックアップ削除（30日以上前）
find $BACKUP_DIR -name "full_backup_*.sql.gz" -mtime +30 -delete
```

### 8.2 定期メンテナンス

```sql
-- 定期実行推奨コマンド
VACUUM ANALYZE reservations;
REINDEX INDEX CONCURRENTLY idx_reservations_staff_start_at;

-- 古いログデータ削除（90日以上前）
DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '90 days';
DELETE FROM notification_logs WHERE created_at < NOW() - INTERVAL '30 days';
```

## 9. Go/GORM との連携考慮事項

### 9.1 GORM タグ例

```go
type Reservation struct {
    ID                    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    CustomerID            uuid.UUID `gorm:"type:uuid;not null" json:"customer_id"`
    StaffID               uuid.UUID `gorm:"type:uuid;not null" json:"staff_id"`
    MenuID                uuid.UUID `gorm:"type:uuid;not null" json:"menu_id"`
    StartAt               time.Time `gorm:"type:timestamp with time zone;not null" json:"start_at"`
    EndAt                 time.Time `gorm:"type:timestamp with time zone;not null" json:"end_at"`
    Status                string    `gorm:"type:varchar(20);default:pending" json:"status"`
    TotalPrice            int       `gorm:"not null" json:"total_price"`
    TotalDurationMinutes  int       `gorm:"not null" json:"total_duration_minutes"`
    CreatedAt             time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
    UpdatedAt             time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
    
    // Relations
    Customer Customer `json:"customer,omitempty"`
    Staff    Staff    `json:"staff,omitempty"`
    Menu     Menu     `json:"menu,omitempty"`
    Options  []Option `gorm:"many2many:reservation_options" json:"options,omitempty"`
}
```

## 10. 変更履歴

| 版 | 日付 | 変更内容 | 作成者 |
|----|------|----------|--------|
| v1.0 | 2025-06-29 | 初版作成 | Claude |