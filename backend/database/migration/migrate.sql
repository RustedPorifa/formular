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




-- Индексы
CREATE INDEX IF NOT EXISTS idx_users_mail ON users(email);

-- Создаем новую таблицу variants с полной информацией
CREATE TABLE IF NOT EXISTS variants (
    uuid VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    class VARCHAR(10) NOT NULL,
    subject VARCHAR(50) NOT NULL,
    solved BOOLEAN DEFAULT FALSE,
    pdf_file_path VARCHAR(500) NOT NULL,
    video_file_path VARCHAR(500)
);
