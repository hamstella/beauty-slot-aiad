-- 美容室予約管理アプリ - インデックス・ビュー・関数削除（ロールバック用）

-- ================================
-- 関数削除
-- ================================

DROP FUNCTION IF EXISTS find_available_slots(UUID, DATE, INTEGER);
DROP FUNCTION IF EXISTS find_available_staff(UUID, DATE, INTEGER);
DROP FUNCTION IF EXISTS create_monthly_partitions(TEXT, DATE, INTEGER);

-- ================================
-- ビュー削除
-- ================================

DROP VIEW IF EXISTS reservation_details;
DROP VIEW IF EXISTS staff_availability;
DROP VIEW IF EXISTS staff_with_labels;
DROP VIEW IF EXISTS menus_with_labels;

-- ================================
-- インデックス削除
-- ================================

-- ログテーブルインデックス
DROP INDEX IF EXISTS idx_notification_logs_recipient;
DROP INDEX IF EXISTS idx_notification_logs_channel;
DROP INDEX IF EXISTS idx_notification_logs_sent_at;
DROP INDEX IF EXISTS idx_notification_logs_status;
DROP INDEX IF EXISTS idx_notification_logs_reservation_id;

DROP INDEX IF EXISTS idx_audit_logs_new_values;
DROP INDEX IF EXISTS idx_audit_logs_old_values;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_entity;

-- 予約関連インデックス
DROP INDEX IF EXISTS idx_reservation_options_option_id;
DROP INDEX IF EXISTS idx_reservation_options_reservation_id;

DROP INDEX IF EXISTS idx_reservations_no_overlap;
DROP INDEX IF EXISTS idx_reservations_staff_status;
DROP INDEX IF EXISTS idx_reservations_start_at;
DROP INDEX IF EXISTS idx_reservations_status;
DROP INDEX IF EXISTS idx_reservations_customer_id;
DROP INDEX IF EXISTS idx_reservations_staff_start_at;

-- シフトインデックス
DROP INDEX IF EXISTS idx_shifts_staff_id;
DROP INDEX IF EXISTS idx_shifts_work_date;
DROP INDEX IF EXISTS idx_shifts_staff_date;

-- 関連テーブルインデックス
DROP INDEX IF EXISTS idx_menu_labels_label_id;
DROP INDEX IF EXISTS idx_menu_labels_menu_id;

DROP INDEX IF EXISTS idx_staff_labels_label_id;
DROP INDEX IF EXISTS idx_staff_labels_staff_id;

-- 基本テーブルインデックス
DROP INDEX IF EXISTS idx_labels_name;

DROP INDEX IF EXISTS idx_options_name_fulltext;
DROP INDEX IF EXISTS idx_options_add_price;
DROP INDEX IF EXISTS idx_options_is_active;

DROP INDEX IF EXISTS idx_menus_name_fulltext;
DROP INDEX IF EXISTS idx_menus_duration;
DROP INDEX IF EXISTS idx_menus_price;
DROP INDEX IF EXISTS idx_menus_is_active;

DROP INDEX IF EXISTS idx_staff_name;
DROP INDEX IF EXISTS idx_staff_is_active;

DROP INDEX IF EXISTS idx_customers_created_at;
DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_customers_phone;