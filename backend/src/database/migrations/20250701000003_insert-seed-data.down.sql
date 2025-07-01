-- 美容室予約管理アプリ - サンプルデータ削除（ロールバック用）

-- ================================
-- 管理用関数削除
-- ================================
DROP FUNCTION IF EXISTS bulk_insert_shifts(UUID, DATE, DATE, TIME, TIME, INTEGER[]);
DROP FUNCTION IF EXISTS cleanup_old_data();

-- ================================
-- サンプルデータ削除（外部キー制約順）
-- ================================

-- 予約関連データ削除
DELETE FROM reservation_options;
DELETE FROM reservations;

-- シフトデータ削除
DELETE FROM shifts;

-- 顧客データ削除
DELETE FROM customers;

-- オプションデータ削除
DELETE FROM options;

-- メニュー関連データ削除
DELETE FROM menu_labels;
DELETE FROM menus;

-- スタッフ関連データ削除
DELETE FROM staff_labels;
DELETE FROM staff;

-- ラベルデータ削除
DELETE FROM labels;