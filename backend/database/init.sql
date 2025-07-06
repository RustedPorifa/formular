-- Создание базы данных
CREATE DATABASE learning_platform;

-- Подключение к БД
\c learning_platform

-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица учебных материалов
CREATE TABLE materials (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(10) CHECK (type IN ('pdf', 'video', 'link')),
    file_path VARCHAR(255),
    url VARCHAR(255),
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для ускорения поиска
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_materials_type ON materials(type);