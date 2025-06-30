-- 美容室予約管理アプリ データベース初期化SQL
-- 実行順序: 001_create_tables.sql → 002_create_indexes.sql → 003_insert_seed_data.sql

-- データベース設定
SET timezone = 'Asia/Tokyo';

-- UUID拡張有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 更新日時自動更新関数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ================================
-- 基本マスターテーブル
-- ================================

-- 顧客テーブル
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

-- スタッフテーブル
CREATE TABLE staff (
    staff_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT staff_name_check CHECK (LENGTH(TRIM(name)) > 0)
);

-- メニューテーブル
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

-- オプションテーブル
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

-- ラベルテーブル
CREATE TABLE labels (
    label_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT labels_name_check CHECK (LENGTH(TRIM(name)) > 0)
);

-- ================================
-- 関連テーブル
-- ================================

-- スタッフラベル関連テーブル
CREATE TABLE staff_labels (
    staff_id UUID NOT NULL,
    label_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (staff_id, label_id),
    FOREIGN KEY (staff_id) REFERENCES staff(staff_id) ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES labels(label_id) ON DELETE CASCADE
);

-- メニューラベル関連テーブル
CREATE TABLE menu_labels (
    menu_id UUID NOT NULL,
    label_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (menu_id, label_id),
    FOREIGN KEY (menu_id) REFERENCES menus(menu_id) ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES labels(label_id) ON DELETE CASCADE
);

-- ================================
-- シフト・予約テーブル
-- ================================

-- シフトテーブル
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

-- 予約テーブル
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

-- 予約オプション関連テーブル
CREATE TABLE reservation_options (
    reservation_id UUID NOT NULL,
    option_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (reservation_id, option_id),
    FOREIGN KEY (reservation_id) REFERENCES reservations(reservation_id) ON DELETE CASCADE,
    FOREIGN KEY (option_id) REFERENCES options(option_id) ON DELETE RESTRICT
);

-- ================================
-- ログテーブル（パーティション対応）
-- ================================

-- 監査ログテーブル
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

-- 通知ログテーブル
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

-- ================================
-- パーティション作成（2025年）
-- ================================

-- 監査ログパーティション
CREATE TABLE audit_logs_2025_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE audit_logs_2025_02 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
CREATE TABLE audit_logs_2025_03 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');
CREATE TABLE audit_logs_2025_04 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-04-01') TO ('2025-05-01');
CREATE TABLE audit_logs_2025_05 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-05-01') TO ('2025-06-01');
CREATE TABLE audit_logs_2025_06 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-06-01') TO ('2025-07-01');
CREATE TABLE audit_logs_2025_07 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');
CREATE TABLE audit_logs_2025_08 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');
CREATE TABLE audit_logs_2025_09 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');
CREATE TABLE audit_logs_2025_10 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');
CREATE TABLE audit_logs_2025_11 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
CREATE TABLE audit_logs_2025_12 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

-- 通知ログパーティション
CREATE TABLE notification_logs_2025_01 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE notification_logs_2025_02 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
CREATE TABLE notification_logs_2025_03 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');
CREATE TABLE notification_logs_2025_04 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-04-01') TO ('2025-05-01');
CREATE TABLE notification_logs_2025_05 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-05-01') TO ('2025-06-01');
CREATE TABLE notification_logs_2025_06 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-06-01') TO ('2025-07-01');
CREATE TABLE notification_logs_2025_07 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');
CREATE TABLE notification_logs_2025_08 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');
CREATE TABLE notification_logs_2025_09 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');
CREATE TABLE notification_logs_2025_10 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');
CREATE TABLE notification_logs_2025_11 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
CREATE TABLE notification_logs_2025_12 PARTITION OF notification_logs
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

-- ================================
-- 更新日時自動更新トリガー
-- ================================

CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_staff_updated_at BEFORE UPDATE ON staff
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menus_updated_at BEFORE UPDATE ON menus
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_options_updated_at BEFORE UPDATE ON options
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_labels_updated_at BEFORE UPDATE ON labels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shifts_updated_at BEFORE UPDATE ON shifts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント
COMMENT ON TABLE customers IS '顧客情報';
COMMENT ON TABLE staff IS 'スタッフ情報';
COMMENT ON TABLE menus IS 'メニュー情報';
COMMENT ON TABLE options IS 'オプション情報';
COMMENT ON TABLE labels IS 'ラベル情報（メニュー・スタッフの分類用）';
COMMENT ON TABLE staff_labels IS 'スタッフ対応ラベル関連';
COMMENT ON TABLE menu_labels IS 'メニューラベル関連';
COMMENT ON TABLE shifts IS 'スタッフシフト情報';
COMMENT ON TABLE reservations IS '予約情報';
COMMENT ON TABLE reservation_options IS '予約オプション関連';
COMMENT ON TABLE audit_logs IS '監査ログ（パーティション対応）';
COMMENT ON TABLE notification_logs IS '通知ログ（パーティション対応）';