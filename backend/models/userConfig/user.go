// Структуры для пользователя
package user

// Обычный профиль пользователя
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Логин пользователя
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWT токен пользователя
type JwtToken struct {
	Token string `json:"token"`
}

// Информация о пользователе для ДБ
type UserInfo struct {
	ID    string
	Name  string
	Email string
	Role  string
}
