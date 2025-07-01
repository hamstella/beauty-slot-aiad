-- 美容室管理システムトークンテーブル（認証・セッション管理用）
CREATE TABLE tokens(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    token           VARCHAR(255)    NOT NULL UNIQUE,
    user_id         UUID            NOT NULL,
    type            VARCHAR(255)    NOT NULL CHECK (type IN ('access', 'refresh', 'reset_password', 'email_verification')),
    expires         TIMESTAMP       NOT NULL,
    is_revoked      BOOLEAN         DEFAULT FALSE NOT NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- インデックス追加
CREATE INDEX idx_tokens_token ON tokens(token);
CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_type ON tokens(type);
CREATE INDEX idx_tokens_expires ON tokens(expires);
CREATE INDEX idx_tokens_is_revoked ON tokens(is_revoked);

-- 更新日時自動更新トリガー
CREATE TRIGGER update_tokens_updated_at BEFORE UPDATE ON tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント
COMMENT ON TABLE tokens IS '美容室システム認証トークン（JWT、リフレッシュトークン等）';
COMMENT ON COLUMN tokens.type IS 'トークン種別（access, refresh, reset_password, email_verification）';
COMMENT ON COLUMN tokens.is_revoked IS 'トークン無効化フラグ';
