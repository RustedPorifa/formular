-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    mail VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL
);

-- Таблица вариантов
CREATE TABLE IF NOT EXISTS variants (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица связи пользователь-вариант
CREATE TABLE IF NOT EXISTS user_completed_variants (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    variant_id INT NOT NULL REFERENCES variants(id) ON DELETE CASCADE,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, variant_id)
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_users_mail ON users(mail);
CREATE INDEX IF NOT EXISTS idx_completed_variants_user ON user_completed_variants(user_id);