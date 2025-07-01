-- 美容室管理システムユーザーテーブル（管理者・スタッフ認証用）
CREATE TABLE users(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(255)    NOT NULL,
    email           VARCHAR(255)    NOT NULL UNIQUE,
    password        VARCHAR(255)    NOT NULL,
    role            VARCHAR(255)    NOT NULL CHECK (role IN ('admin', 'staff', 'customer')),
    staff_id        UUID            REFERENCES staff(staff_id) ON DELETE SET NULL,
    customer_id     UUID            REFERENCES customers(customer_id) ON DELETE SET NULL,
    verified_email  BOOLEAN         DEFAULT FALSE  NOT NULL,
    is_active       BOOLEAN         DEFAULT TRUE   NOT NULL,
    last_login_at   TIMESTAMP       NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    
    CONSTRAINT users_role_staff_check CHECK (
        (role = 'staff' AND staff_id IS NOT NULL AND customer_id IS NULL) OR
        (role = 'customer' AND customer_id IS NOT NULL AND staff_id IS NULL) OR
        (role = 'admin' AND staff_id IS NULL AND customer_id IS NULL)
    )
);

-- インデックス追加
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_staff_id ON users(staff_id);
CREATE INDEX idx_users_customer_id ON users(customer_id);
CREATE INDEX idx_users_is_active ON users(is_active);

-- 更新日時自動更新トリガー
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント
COMMENT ON TABLE users IS '美容室システムユーザー（管理者・スタッフ・顧客認証用）';
COMMENT ON COLUMN users.staff_id IS 'スタッフの場合のスタッフID（外部キー）';
COMMENT ON COLUMN users.customer_id IS '顧客の場合の顧客ID（外部キー）';
