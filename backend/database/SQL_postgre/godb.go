package godb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	user "formular/backend/models/userConfig"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
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
	migrationSQL, err := os.ReadFile("backend/database/migration/migrate.sql")
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

func AddUser(ctx context.Context, user *user.User) error {
	sql := `
        INSERT INTO users (id, name, email, password, role, purchased_grades, is_authenticated) 
        VALUES (@id, @name, @email, @password, @role, @purchased_grades, @is_authenticated)
    `
	// Преобразуем nil в пустой массив
	purchasedGrades := user.PurchasedGrades
	if purchasedGrades == nil {
		purchasedGrades = []string{}
	}
	args := pgx.NamedArgs{
		"id":               user.ID,
		"name":             user.Name,
		"email":            user.Email,
		"password":         user.Password,
		"role":             user.Role,
		"purchased_grades": purchasedGrades,
		"is_authenticated": user.IsAuthenticated, // Добавляем это поле
	}
	_, err := pool.Exec(ctx, sql, args)
	return err
}

// AddAdmin создает администратора с указанными данными
func AddAdmin() error {
	name := "admin"
	email := "formulyarka@yandex.ru"
	password := "$2a$14$9I7lYldd867nz/Oe4hlhYeEI8nM/xTZbviS5CIBEvyP6cCweG9BzK"
	admin := &user.User{
		ID:              uuid.New().String(),
		Name:            name,
		Email:           email,
		Password:        password,
		Role:            "Admin",
		IsAuthenticated: true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	println("making...")
	return AddUser(ctx, admin)
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
        SELECT id, name, email, password, role, purchased_grades 
        FROM users 
        WHERE email = @email
    `

	err := pool.QueryRow(ctx, query, pgx.NamedArgs{"email": credentials.Email}).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.PurchasedGrades,
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
		SELECT id, name, email, password, role, is_authenticated, purchased_grades
		FROM users WHERE email = @email
	`
	user := &user.User{}
	err := pool.QueryRow(ctx, sql, pgx.NamedArgs{"email": email}).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsAuthenticated,
		&user.PurchasedGrades,
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
func GetUserInfoByID(ctx context.Context, userID string) (*user.User, error) {
	query := `
		SELECT 
			id,
			name,
			email,
			role,
			purchased_grades  -- Берем данные напрямую из таблицы users
		FROM users
		WHERE id = @userID
	`

	var u user.User
	err := pool.QueryRow(ctx, query, pgx.NamedArgs{"userID": userID}).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Role,
		&u.PurchasedGrades, // Сканируем напрямую в поле структуры
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Заполняем дополнительные поля
	u.Password = ""          // Пароль не возвращается из соображений безопасности
	u.IsAuthenticated = true // Пользователь аутентифицирован (так как найден по ID)

	return &u, nil
}

// CreateAnonymousUser создает временного анонимного пользователя
func CreateAnonymousUser(ctx context.Context) (*user.User, error) {
	id := uuid.New().String()
	anonEmail := "anon_" + id + "@example.com"

	sql := `
        INSERT INTO users (id, name, email, password, role, is_authenticated)
        VALUES (@id, 'Anonymous', @email, NULL, 'anonymous', false)
        RETURNING id, name, email, role, is_authenticated, purchased_grades
    `

	u := &user.User{}
	err := pool.QueryRow(ctx, sql,
		pgx.NamedArgs{
			"id":    id,
			"email": anonEmail,
		}).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Role,
		&u.IsAuthenticated,
		&u.PurchasedGrades,
	)

	return u, err
}

// SetUserAuthenticatedAndRole обновляет статус аутентификации и устанавливает роль Member для пользователя
func SetUserAuthenticatedAndRole(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET 
			is_authenticated = true,
			role = 'Member'
		WHERE id = $1
	`

	result, err := pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update user authentication status and role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

// DeleteAllUnauthenticatedUsers удаляет всех неавторизованных пользователей и их связанные данные
func DeleteAllUnauthenticatedUsers(ctx context.Context) (int64, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Всегда пытаемся откатить при ошибке

	deleteVariantsSQL := `
        DELETE FROM user_completed_variants
        WHERE user_id IN (
            SELECT id FROM users WHERE is_authenticated = false
        )
    `
	_, err = tx.Exec(ctx, deleteVariantsSQL)
	if err != nil {
		return 0, fmt.Errorf("failed to delete related variants: %w", err)
	}

	deleteUsersSQL := `
        DELETE FROM users
        WHERE is_authenticated = false
    `
	tag, err := tx.Exec(ctx, deleteUsersSQL)
	if err != nil {
		return 0, fmt.Errorf("failed to delete users: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tag.RowsAffected(), nil
}

// FindUsersByPurchasedGrade возвращает пользователей, купивших указанный класс
func FindUsersByPurchasedGrade(ctx context.Context, grade string) ([]user.User, error) {
	query := `
        SELECT 
            id, name, email, password, role, is_authenticated, purchased_grades
        FROM users
        WHERE @grade = ANY(purchased_grades)
    `

	rows, err := pool.Query(ctx, query, pgx.NamedArgs{"grade": grade})
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Password,
			&u.Role,
			&u.IsAuthenticated,
			&u.PurchasedGrades,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

// GetPurchasedGradesByUserID возвращает массив купленных классов для указанного пользователя
func GetPurchasedGradesByUserID(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT purchased_grades FROM users WHERE id = @userID`

	var grades []string
	err := pool.QueryRow(ctx, query, pgx.NamedArgs{"userID": userID}).Scan(&grades)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get purchased grades: %w", err)
	}

	return grades, nil
}
