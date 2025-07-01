-- インデックス削除
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_customer_id;
DROP INDEX IF EXISTS idx_users_staff_id;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;

-- トリガー削除
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- テーブル削除
DROP TABLE IF EXISTS users;