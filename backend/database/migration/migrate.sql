-- Разрешаем NULL для пароля
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;

-- Добавляем флаг аутентификации
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS is_authenticated BOOLEAN NOT NULL DEFAULT FALSE;

-- Таблица пользователей (актуальная версия)
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    role VARCHAR(50) NOT NULL,
    is_authenticated BOOLEAN NOT NULL DEFAULT FALSE
);

-- Остальные таблицы без изменений
CREATE TABLE IF NOT EXISTS variants (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_completed_variants (
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    variant_id INT NOT NULL REFERENCES variants(id) ON DELETE CASCADE,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, variant_id)
);

CREATE INDEX IF NOT EXISTS idx_users_mail ON users(email);
CREATE INDEX IF NOT EXISTS idx_completed_variants_user ON user_completed_variants(user_id);