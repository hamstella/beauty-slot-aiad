-- 美容室予約管理アプリ - 初期データ・サンプルデータ投入
-- Generated from: database/migrations/003_insert_seed_data.sql

-- ================================
-- ラベルデータ
-- ================================
INSERT INTO labels (label_id, name) VALUES
('11111111-1111-1111-1111-111111111111', 'カット'),
('22222222-2222-2222-2222-222222222222', 'カラー'),
('33333333-3333-3333-3333-333333333333', 'パーマ'),
('44444444-4444-4444-4444-444444444444', 'トリートメント'),
('55555555-5555-5555-5555-555555555555', 'メンズカット'),
('66666666-6666-6666-6666-666666666666', 'ブリーチ'),
('77777777-7777-7777-7777-777777777777', 'ヘッドスパ'),
('88888888-8888-8888-8888-888888888888', 'セット・アップ');

-- ================================
-- スタッフデータ
-- ================================
INSERT INTO staff (staff_id, name, is_active) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'タナカ ミカ', true),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'サトウ タロウ', true),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'ヤマダ ハナコ', true),
('dddddddd-dddd-dddd-dddd-dddddddddddd', 'スズキ ジロウ', true);

-- ================================
-- スタッフ対応ラベル設定
-- ================================
-- 田中 美香：カット、カラー、トリートメント
INSERT INTO staff_labels (staff_id, label_id) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '22222222-2222-2222-2222-222222222222'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '44444444-4444-4444-4444-444444444444');

-- 佐藤 太郎：メンズカット、カット、ブリーチ
INSERT INTO staff_labels (staff_id, label_id) VALUES
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '55555555-5555-5555-5555-555555555555'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '66666666-6666-6666-6666-666666666666');

-- 山田 花子：カット、パーマ、ヘッドスパ、セット・アップ
INSERT INTO staff_labels (staff_id, label_id) VALUES
('cccccccc-cccc-cccc-cccc-cccccccccccc', '11111111-1111-1111-1111-111111111111'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', '33333333-3333-3333-3333-333333333333'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', '77777777-7777-7777-7777-777777777777'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', '88888888-8888-8888-8888-888888888888');

-- 鈴木 次郎：カット、カラー、ブリーチ、トリートメント
INSERT INTO staff_labels (staff_id, label_id) VALUES
('dddddddd-dddd-dddd-dddd-dddddddddddd', '11111111-1111-1111-1111-111111111111'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', '22222222-2222-2222-2222-222222222222'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', '66666666-6666-6666-6666-666666666666'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', '44444444-4444-4444-4444-444444444444');

-- ================================
-- メニューデータ
-- ================================
INSERT INTO menus (menu_id, name, duration_minutes, price, is_active) VALUES
('m1111111-1111-1111-1111-111111111111', 'レディースカット', 60, 4500, true),
('m2222222-2222-2222-2222-222222222222', 'メンズカット', 45, 3500, true),
('m3333333-3333-3333-3333-333333333333', 'カラーリング（ルート）', 90, 6500, true),
('m4444444-4444-4444-4444-444444444444', 'カラーリング（フル）', 120, 8500, true),
('m5555555-5555-5555-5555-555555555555', 'パーマ', 150, 12000, true),
('m6666666-6666-6666-6666-666666666666', 'ブリーチ', 180, 15000, true),
('m7777777-7777-7777-7777-777777777777', 'ヘッドスパ', 30, 2500, true),
('m8888888-8888-8888-8888-888888888888', 'セット・アップ', 45, 3000, true);

-- ================================
-- メニューラベル設定
-- ================================
-- レディースカット
INSERT INTO menu_labels (menu_id, label_id) VALUES
('m1111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111');

-- メンズカット
INSERT INTO menu_labels (menu_id, label_id) VALUES
('m2222222-2222-2222-2222-222222222222', '55555555-5555-5555-5555-555555555555');

-- カラーリング（ルート）
INSERT INTO menu_labels (menu_id, label_id) VALUES
('m3333333-3333-3333-3333-333333333333', '22222222-2222-2222-2222-222222222222');

-- カラーリング（フル）
INSERT INTO menu_labels (menu_id, label_id) VALUES
('m4444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222');

-- パーマ
INSERT INTO menu_labels (menu_id, label_id) VALUES
('m5555555-5555-5555-5555-555555555555', '33333333-3333-3333-3333-333333333333');

-- ブリーチ
INSERT INTO menu_labels (menu_id, label_id) VALUES
('m6666666-6666-6666-6666-666666666666', '66666666-6666-6666-6666-666666666666');

-- ヘッドスパ
INSERT INTO menu_labels (menu_id, label_id) VALUES
('m7777777-7777-7777-7777-777777777777', '77777777-7777-7777-7777-777777777777');

-- セット・アップ
INSERT INTO menu_labels (menu_id, label_id) VALUES
('m8888888-8888-8888-8888-888888888888', '88888888-8888-8888-8888-888888888888');

-- ================================
-- オプションデータ
-- ================================
INSERT INTO options (option_id, name, add_duration_minutes, add_price, is_active) VALUES
('o1111111-1111-1111-1111-111111111111', 'トリートメント追加', 15, 1500, true),
('o2222222-2222-2222-2222-222222222222', 'ヘッドマッサージ', 10, 1000, true),
('o3333333-3333-3333-3333-333333333333', '眉カット', 15, 800, true),
('o4444444-4444-4444-4444-444444444444', 'シャンプー・ブロー', 20, 1200, true),
('o5555555-5555-5555-5555-555555555555', 'スタイリング', 15, 1000, true);

-- ================================
-- 顧客データ（サンプル）
-- ================================
INSERT INTO customers (customer_id, name, phone, email) VALUES
('c1111111-1111-1111-1111-111111111111', 'ヤマダ タロウ', '090-1234-5678', 'yamada.taro@example.com'),
('c2222222-2222-2222-2222-222222222222', 'タナカ ハナコ', '080-9876-5432', 'tanaka.hanako@example.com'),
('c3333333-3333-3333-3333-333333333333', 'サトウ ジロウ', '070-1111-2222', 'sato.jiro@example.com'),
('c4444444-4444-4444-4444-444444444444', 'スズキ ミカ', '090-3333-4444', 'suzuki.mika@example.com'),
('c5555555-5555-5555-5555-555555555555', 'タカハシ ヒロシ', '080-5555-6666', 'takahashi.hiroshi@example.com');

-- ================================
-- シフトデータ（今後7日間）
-- ================================
-- 田中 美香のシフト（月火水金土）
INSERT INTO shifts (staff_id, work_date, start_time, end_time) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', CURRENT_DATE + 1, '09:00', '18:00'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', CURRENT_DATE + 2, '09:00', '18:00'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', CURRENT_DATE + 3, '09:00', '18:00'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', CURRENT_DATE + 5, '09:00', '18:00'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', CURRENT_DATE + 6, '09:00', '17:00');

-- 佐藤 太郎のシフト（火水木金土）
INSERT INTO shifts (staff_id, work_date, start_time, end_time) VALUES
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', CURRENT_DATE + 2, '10:00', '19:00'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', CURRENT_DATE + 3, '10:00', '19:00'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', CURRENT_DATE + 4, '10:00', '19:00'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', CURRENT_DATE + 5, '10:00', '19:00'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', CURRENT_DATE + 6, '10:00', '18:00');

-- 山田 花子のシフト（月水木土日）
INSERT INTO shifts (staff_id, work_date, start_time, end_time) VALUES
('cccccccc-cccc-cccc-cccc-cccccccccccc', CURRENT_DATE + 1, '09:30', '17:30'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', CURRENT_DATE + 3, '09:30', '17:30'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', CURRENT_DATE + 4, '09:30', '17:30'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', CURRENT_DATE + 6, '09:30', '17:30'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', CURRENT_DATE + 7, '09:30', '16:30');

-- 鈴木 次郎のシフト（月火木金土）
INSERT INTO shifts (staff_id, work_date, start_time, end_time) VALUES
('dddddddd-dddd-dddd-dddd-dddddddddddd', CURRENT_DATE + 1, '10:00', '18:00'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', CURRENT_DATE + 2, '10:00', '18:00'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', CURRENT_DATE + 4, '10:00', '18:00'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', CURRENT_DATE + 5, '10:00', '18:00'),
('dddddddd-dddd-dddd-dddd-dddddddddddd', CURRENT_DATE + 6, '10:00', '17:00');

-- ================================
-- サンプル予約データ
-- ================================
-- 明日の予約
INSERT INTO reservations (reservation_id, customer_id, staff_id, menu_id, start_at, end_at, status, total_price, total_duration_minutes) VALUES
('r1111111-1111-1111-1111-111111111111', 'c1111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'm1111111-1111-1111-1111-111111111111', 
 (CURRENT_DATE + 1) + TIME '10:00', (CURRENT_DATE + 1) + TIME '11:15', 'confirmed', 6000, 75),
 
('r2222222-2222-2222-2222-222222222222', 'c2222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'm2222222-2222-2222-2222-222222222222', 
 (CURRENT_DATE + 2) + TIME '14:00', (CURRENT_DATE + 2) + TIME '14:45', 'confirmed', 3500, 45);

-- 予約オプション設定
INSERT INTO reservation_options (reservation_id, option_id) VALUES
('r1111111-1111-1111-1111-111111111111', 'o1111111-1111-1111-1111-111111111111'); -- レディースカット + トリートメント

-- ================================
-- 管理用関数・プロシージャ
-- ================================

-- シフト一括登録関数
CREATE OR REPLACE FUNCTION bulk_insert_shifts(
    p_staff_id UUID,
    p_start_date DATE,
    p_end_date DATE,
    p_start_time TIME,
    p_end_time TIME,
    p_work_days INTEGER[] -- 曜日配列（0=日曜、1=月曜...6=土曜）
)
RETURNS INTEGER AS $$
DECLARE
    current_date DATE;
    insert_count INTEGER := 0;
BEGIN
    current_date := p_start_date;
    
    WHILE current_date <= p_end_date LOOP
        -- 指定曜日の場合にシフトを登録
        IF EXTRACT(DOW FROM current_date)::INTEGER = ANY(p_work_days) THEN
            INSERT INTO shifts (staff_id, work_date, start_time, end_time)
            VALUES (p_staff_id, current_date, p_start_time, p_end_time)
            ON CONFLICT (staff_id, work_date) DO NOTHING;
            
            IF FOUND THEN
                insert_count := insert_count + 1;
            END IF;
        END IF;
        
        current_date := current_date + 1;
    END LOOP;
    
    RETURN insert_count;
END;
$$ LANGUAGE plpgsql;

-- データクリーンアップ関数
CREATE OR REPLACE FUNCTION cleanup_old_data()
RETURNS TEXT AS $$
DECLARE
    deleted_audit_logs INTEGER;
    deleted_notification_logs INTEGER;
    result_message TEXT;
BEGIN
    -- 90日以上前の監査ログ削除
    DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '90 days';
    GET DIAGNOSTICS deleted_audit_logs = ROW_COUNT;
    
    -- 30日以上前の通知ログ削除
    DELETE FROM notification_logs WHERE created_at < NOW() - INTERVAL '30 days';
    GET DIAGNOSTICS deleted_notification_logs = ROW_COUNT;
    
    result_message := format('Deleted: %s audit logs, %s notification logs', 
                           deleted_audit_logs, deleted_notification_logs);
    
    RETURN result_message;
END;
$$ LANGUAGE plpgsql;

-- ================================
-- コメント
-- ================================
COMMENT ON FUNCTION bulk_insert_shifts(UUID, DATE, DATE, TIME, TIME, INTEGER[]) IS 'スタッフシフト一括登録関数';
COMMENT ON FUNCTION cleanup_old_data() IS '古いログデータクリーンアップ関数';

-- ================================
-- 統計情報更新
-- ================================
ANALYZE customers;
ANALYZE staff;
ANALYZE menus;
ANALYZE options;
ANALYZE labels;
ANALYZE staff_labels;
ANALYZE menu_labels;
ANALYZE shifts;
ANALYZE reservations;
ANALYZE reservation_options;