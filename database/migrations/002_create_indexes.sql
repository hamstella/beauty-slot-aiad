-- インデックス・制約・ビュー・関数作成SQL
-- 実行順序: 001_create_tables.sql → 002_create_indexes.sql → 003_insert_seed_data.sql

-- ================================
-- 基本インデックス
-- ================================

-- customers テーブル
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_created_at ON customers(created_at);

-- staff テーブル
CREATE INDEX idx_staff_is_active ON staff(is_active);
CREATE INDEX idx_staff_name ON staff(name);

-- menus テーブル
CREATE INDEX idx_menus_is_active ON menus(is_active);
CREATE INDEX idx_menus_price ON menus(price);
CREATE INDEX idx_menus_duration ON menus(duration_minutes);
CREATE INDEX idx_menus_name_fulltext ON menus USING gin(to_tsvector('japanese', name));

-- options テーブル
CREATE INDEX idx_options_is_active ON options(is_active);
CREATE INDEX idx_options_add_price ON options(add_price);
CREATE INDEX idx_options_name_fulltext ON options USING gin(to_tsvector('japanese', name));

-- labels テーブル
CREATE INDEX idx_labels_name ON labels(name);

-- ================================
-- 関連テーブルインデックス
-- ================================

-- staff_labels テーブル
CREATE INDEX idx_staff_labels_staff_id ON staff_labels(staff_id);
CREATE INDEX idx_staff_labels_label_id ON staff_labels(label_id);

-- menu_labels テーブル
CREATE INDEX idx_menu_labels_menu_id ON menu_labels(menu_id);
CREATE INDEX idx_menu_labels_label_id ON menu_labels(label_id);

-- ================================
-- シフト・予約関連インデックス
-- ================================

-- shifts テーブル
CREATE INDEX idx_shifts_staff_date ON shifts(staff_id, work_date);
CREATE INDEX idx_shifts_work_date ON shifts(work_date);
CREATE INDEX idx_shifts_staff_id ON shifts(staff_id);

-- reservations テーブル（重要なクエリパターンに最適化）
CREATE INDEX idx_reservations_staff_start_at ON reservations(staff_id, start_at);
CREATE INDEX idx_reservations_customer_id ON reservations(customer_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_start_at ON reservations(start_at);
CREATE INDEX idx_reservations_staff_status ON reservations(staff_id, status);

-- 重複予約防止制約（同一スタッフ・同一時間帯・非キャンセル）
CREATE UNIQUE INDEX idx_reservations_no_overlap 
ON reservations(staff_id, start_at, end_at) 
WHERE status NOT IN ('cancelled');

-- reservation_options テーブル
CREATE INDEX idx_reservation_options_reservation_id ON reservation_options(reservation_id);
CREATE INDEX idx_reservation_options_option_id ON reservation_options(option_id);

-- ================================
-- ログテーブルインデックス
-- ================================

-- audit_logs テーブル
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_old_values ON audit_logs USING gin(old_values);
CREATE INDEX idx_audit_logs_new_values ON audit_logs USING gin(new_values);

-- notification_logs テーブル
CREATE INDEX idx_notification_logs_reservation_id ON notification_logs(reservation_id);
CREATE INDEX idx_notification_logs_status ON notification_logs(status);
CREATE INDEX idx_notification_logs_sent_at ON notification_logs(sent_at);
CREATE INDEX idx_notification_logs_channel ON notification_logs(channel);
CREATE INDEX idx_notification_logs_recipient ON notification_logs(recipient);

-- ================================
-- ビュー定義
-- ================================

-- 予約詳細ビュー
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
            ) ORDER BY opt.name
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

-- スタッフ空き状況ビュー
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
                'status', r.status,
                'customer_name', c.name
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
LEFT JOIN customers c ON r.customer_id = c.customer_id
WHERE s.is_active = true
GROUP BY s.staff_id, s.name, sh.work_date, sh.start_time, sh.end_time;

-- スタッフ対応ラベル一覧ビュー
CREATE VIEW staff_with_labels AS
SELECT 
    s.staff_id,
    s.name AS staff_name,
    s.is_active,
    COALESCE(
        json_agg(
            json_build_object(
                'label_id', l.label_id,
                'label_name', l.name
            )
            ORDER BY l.name
        ) FILTER (WHERE l.label_id IS NOT NULL),
        '[]'::json
    ) AS labels
FROM staff s
LEFT JOIN staff_labels sl ON s.staff_id = sl.staff_id
LEFT JOIN labels l ON sl.label_id = l.label_id
GROUP BY s.staff_id, s.name, s.is_active;

-- メニュー・ラベル一覧ビュー
CREATE VIEW menus_with_labels AS
SELECT 
    m.menu_id,
    m.name AS menu_name,
    m.duration_minutes,
    m.price,
    m.is_active,
    COALESCE(
        json_agg(
            json_build_object(
                'label_id', l.label_id,
                'label_name', l.name
            )
            ORDER BY l.name
        ) FILTER (WHERE l.label_id IS NOT NULL),
        '[]'::json
    ) AS labels
FROM menus m
LEFT JOIN menu_labels ml ON m.menu_id = ml.menu_id
LEFT JOIN labels l ON ml.label_id = l.label_id
GROUP BY m.menu_id, m.name, m.duration_minutes, m.price, m.is_active;

-- ================================
-- 空き時間検索関数
-- ================================

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
    shift_end_datetime TIMESTAMP WITH TIME ZONE;
    rec RECORD;
BEGIN
    -- スタッフのシフト情報取得
    SELECT start_time, end_time INTO shift_start, shift_end
    FROM shifts 
    WHERE staff_id = p_staff_id AND work_date = p_date;
    
    IF shift_start IS NULL THEN
        RETURN; -- シフトなし
    END IF;
    
    -- 開始時刻・終了時刻を設定
    slot_start := (p_date + shift_start) AT TIME ZONE 'Asia/Tokyo';
    shift_end_datetime := (p_date + shift_end) AT TIME ZONE 'Asia/Tokyo';
    
    -- 予約済み時間を取得してソート
    FOR rec IN 
        SELECT start_at, end_at
        FROM reservations 
        WHERE staff_id = p_staff_id 
            AND DATE(start_at AT TIME ZONE 'Asia/Tokyo') = p_date
            AND status NOT IN ('cancelled')
        ORDER BY start_at
    LOOP
        slot_end := slot_start + (p_duration_minutes || ' minutes')::INTERVAL;
        
        -- 空き時間が予約時間と重複しないかチェック
        IF slot_end <= rec.start_at THEN
            available_start_time := slot_start;
            available_end_time := slot_end;
            RETURN NEXT;
        END IF;
        
        -- 次の開始時刻を予約終了時刻に設定
        slot_start := rec.end_at;
    END LOOP;
    
    -- 最後の予約後の空き時間をチェック
    slot_end := slot_start + (p_duration_minutes || ' minutes')::INTERVAL;
    IF slot_end <= shift_end_datetime THEN
        available_start_time := slot_start;
        available_end_time := slot_end;
        RETURN NEXT;
    END IF;
    
    RETURN;
END;
$$ LANGUAGE plpgsql;

-- ================================
-- スタッフ検索関数（ラベルマッチング）
-- ================================

CREATE OR REPLACE FUNCTION find_available_staff(
    p_menu_id UUID,
    p_date DATE,
    p_duration_minutes INTEGER
)
RETURNS TABLE(
    staff_id UUID,
    staff_name VARCHAR(100)
) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT s.staff_id, s.name
    FROM staff s
    JOIN shifts sh ON s.staff_id = sh.staff_id
    WHERE s.is_active = true
        AND sh.work_date = p_date
        AND s.staff_id IN (
            -- メニューのラベルに対応可能なスタッフ
            SELECT sl.staff_id
            FROM staff_labels sl
            JOIN menu_labels ml ON sl.label_id = ml.label_id
            WHERE ml.menu_id = p_menu_id
        )
        AND EXISTS (
            -- 指定時間の空きがあるスタッフ
            SELECT 1 
            FROM find_available_slots(s.staff_id, p_date, p_duration_minutes)
            LIMIT 1
        )
    ORDER BY s.name;
END;
$$ LANGUAGE plpgsql;

-- ================================
-- パーティション管理関数
-- ================================

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

-- ================================
-- 統計情報・パフォーマンス設定
-- ================================

-- 統計情報の自動更新設定
ALTER TABLE reservations SET (autovacuum_analyze_scale_factor = 0.02);
ALTER TABLE audit_logs SET (autovacuum_analyze_scale_factor = 0.1);
ALTER TABLE notification_logs SET (autovacuum_analyze_scale_factor = 0.1);

-- ================================
-- コメント追加
-- ================================

COMMENT ON VIEW reservation_details IS '予約詳細情報（顧客・スタッフ・メニュー・オプション含む）';
COMMENT ON VIEW staff_availability IS 'スタッフ別空き状況（シフト・予約情報含む）';
COMMENT ON VIEW staff_with_labels IS 'スタッフ対応ラベル一覧';
COMMENT ON VIEW menus_with_labels IS 'メニューラベル一覧';

COMMENT ON FUNCTION find_available_slots(UUID, DATE, INTEGER) IS '指定スタッフ・日付・所要時間での空き時間検索';
COMMENT ON FUNCTION find_available_staff(UUID, DATE, INTEGER) IS 'メニュー対応可能な空きスタッフ検索';
COMMENT ON FUNCTION create_monthly_partitions(TEXT, DATE, INTEGER) IS 'パーティション自動作成関数';