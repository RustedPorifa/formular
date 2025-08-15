-- Снимаем NOT NULL с пароля (если таблица существует)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users'
          AND column_name = 'password'
          AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
    END IF;
EXCEPTION WHEN undefined_table THEN
    RAISE NOTICE 'Table "users" does not exist. Skipping password ALTER.';
END $$;

-- Создаём таблицу users (если не существует) с актуальной структурой
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    role VARCHAR(50) NOT NULL,
    is_authenticated BOOLEAN NOT NULL DEFAULT FALSE,
    purchased_grades TEXT[] NOT NULL DEFAULT '{}' -- Массив купленных классов
);

-- Убедимся, что purchased_grades имеет правильный DEFAULT
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users'
          AND column_name = 'purchased_grades'
          AND (column_default IS NULL OR column_default <> '''{}''::text[]')
    ) THEN
        ALTER TABLE users ALTER COLUMN purchased_grades SET DEFAULT '{}';
    END IF;
EXCEPTION WHEN undefined_table THEN
    RAISE NOTICE 'Table "users" does not exist. Skipping purchased_grades ALTER.';
END $$;

-- Добавляем новые колонки в существующую таблицу (если отсутствуют)
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_authenticated BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS purchased_grades TEXT[] NOT NULL DEFAULT '{}';

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

-- Индексы
CREATE INDEX IF NOT EXISTS idx_users_mail ON users(email);
CREATE INDEX IF NOT EXISTS idx_completed_variants_user ON user_completed_variants(user_id);
