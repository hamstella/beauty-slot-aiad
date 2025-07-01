-- 美容室予約管理アプリ - テーブル削除（ロールバック用）

-- パーティションテーブル削除
DROP TABLE IF EXISTS audit_logs_2025_01;
DROP TABLE IF EXISTS audit_logs_2025_02;
DROP TABLE IF EXISTS audit_logs_2025_03;
DROP TABLE IF EXISTS audit_logs_2025_04;
DROP TABLE IF EXISTS audit_logs_2025_05;
DROP TABLE IF EXISTS audit_logs_2025_06;
DROP TABLE IF EXISTS audit_logs_2025_07;
DROP TABLE IF EXISTS audit_logs_2025_08;
DROP TABLE IF EXISTS audit_logs_2025_09;
DROP TABLE IF EXISTS audit_logs_2025_10;
DROP TABLE IF EXISTS audit_logs_2025_11;
DROP TABLE IF EXISTS audit_logs_2025_12;

DROP TABLE IF EXISTS notification_logs_2025_01;
DROP TABLE IF EXISTS notification_logs_2025_02;
DROP TABLE IF EXISTS notification_logs_2025_03;
DROP TABLE IF EXISTS notification_logs_2025_04;
DROP TABLE IF EXISTS notification_logs_2025_05;
DROP TABLE IF EXISTS notification_logs_2025_06;
DROP TABLE IF EXISTS notification_logs_2025_07;
DROP TABLE IF EXISTS notification_logs_2025_08;
DROP TABLE IF EXISTS notification_logs_2025_09;
DROP TABLE IF EXISTS notification_logs_2025_10;
DROP TABLE IF EXISTS notification_logs_2025_11;
DROP TABLE IF EXISTS notification_logs_2025_12;

-- 関連テーブル削除（外部キー制約順）
DROP TABLE IF EXISTS reservation_options;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS shifts;
DROP TABLE IF EXISTS menu_labels;
DROP TABLE IF EXISTS staff_labels;
DROP TABLE IF EXISTS notification_logs;
DROP TABLE IF EXISTS audit_logs;

-- マスターテーブル削除
DROP TABLE IF EXISTS labels;
DROP TABLE IF EXISTS options;
DROP TABLE IF EXISTS menus;
DROP TABLE IF EXISTS staff;
DROP TABLE IF EXISTS customers;

-- 関数削除
DROP FUNCTION IF EXISTS update_updated_at_column();