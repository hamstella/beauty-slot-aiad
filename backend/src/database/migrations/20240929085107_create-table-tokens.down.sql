-- インデックス削除
DROP INDEX IF EXISTS idx_tokens_is_revoked;
DROP INDEX IF EXISTS idx_tokens_expires;
DROP INDEX IF EXISTS idx_tokens_type;
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_token;

-- トリガー削除
DROP TRIGGER IF EXISTS update_tokens_updated_at ON tokens;

-- テーブル削除
DROP TABLE IF EXISTS tokens;