package godb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	user "formular/backend/models/userConfig"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var pool *pgxpool.Pool

// InitDB инициализирует подключение к БД и выполняет миграции
func InitDB() error {
	// Формирование строки подключения из переменных окружения
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		"?sslmode=disable",
	)

	// Контекст с таймаутом для инициализации
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создание пула соединений
	var err error
	pool, err = pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Выполнение миграций
	if err := runMigrations(ctx); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}

	return nil
}

// Close освобождает ресурсы БД
func Close() {
	if pool != nil {
		pool.Close()
	}
}

// runMigrations выполняет SQL-миграции
func runMigrations(ctx context.Context) error {
	// Чтение файла миграции
	migrationSQL, err := os.ReadFile("backend/database/migrate.sql")
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}

	// Выполнение миграции
	_, err = pool.Exec(ctx, string(migrationSQL))
	if err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}

	return nil
}

// AddUser добавляет нового пользователя
func AddUser(ctx context.Context, user *user.User) error {
	sql := `
        INSERT INTO users (id, name, email, password, role) 
        VALUES (@id, @name, @email, @password, @role)
    `

	args := pgx.NamedArgs{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
		"role":     user.Role,
	}

	_, err := pool.Exec(ctx, sql, args)
	return err
}

// GetUserRole возвращает роль пользователя по email
func GetUserRole(ctx context.Context, email string) (string, error) {
	sql := `SELECT role FROM users WHERE email = @email`
	var role string

	err := pool.QueryRow(ctx, sql, pgx.NamedArgs{"email": email}).Scan(&role)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errors.New("user not found")
	}

	return role, err
}

// Возвращает true и структуру пользователя при успехе, false и пустую структуру при неудаче
func VerifyCredentials(ctx context.Context, credentials user.Credentials) (bool, user.User) {
	// Запрашиваем только необходимые поля пользователя
	var user user.User
	query := `
        SELECT id, name, email, password, role 
        FROM users 
        WHERE email = @email
    `

	err := pool.QueryRow(ctx, query, pgx.NamedArgs{"email": credentials.Email}).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		// Пользователь не найден или ошибка запроса
		return false, user
	}

	// Сравниваем хешированный пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		// Пароль не совпадает
		return false, user
	}

	// Возвращаем true и структуру пользователя
	return true, user
}

// GetUserByEmail возвращает пользователя по email
func GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	sql := `
		SELECT id, name, email, password, role 
		FROM users WHERE email = @email
	`
	user := &user.User{}
	err := pool.QueryRow(ctx, sql, pgx.NamedArgs{"email": email}).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UserExists проверяет существует ли пользователь с указанным ID
func UserExists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = @userID)`
	var exists bool
	err := pool.QueryRow(ctx, query, pgx.NamedArgs{"userID": userID}).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// GetUserInfo возвращает информацию о пользователе по email
func GetUserInfo(ctx context.Context, email string) (*user.UserInfo, error) {
	query := `
        SELECT 
            u.name,
            u.email,
            u.role,
            COUNT(ucv.variant_id) AS completed_count
        FROM users u
        LEFT JOIN user_completed_variants ucv ON u.id = ucv.user_id
        WHERE u.email = @email
        GROUP BY u.id, u.name, u.email, u.role
    `

	var info user.UserInfo
	err := pool.QueryRow(ctx, query, pgx.NamedArgs{"email": email}).Scan(
		&info.ID,
		&info.Name,
		&info.Email,
		&info.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &info, nil
}

// GetUserInfoByID возвращает информацию о пользователе по его ID
func GetUserInfoByID(ctx context.Context, userID string) (*user.UserInfo, error) {
	query := `
        SELECT 
            u.id,
            u.name,
            u.email,
            u.role,
            COUNT(ucv.variant_id) AS completed_count
        FROM users u
        LEFT JOIN user_completed_variants ucv ON u.id = ucv.user_id
        WHERE u.id = @userID
        GROUP BY u.id, u.name, u.email, u.role
    `

	var info user.UserInfo
	err := pool.QueryRow(ctx, query, pgx.NamedArgs{"userID": userID}).Scan(
		&info.ID,
		&info.Name,
		&info.Email,
		&info.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &info, nil
}
