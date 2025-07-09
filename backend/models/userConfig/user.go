package user

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Completed []int
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type UserInfo struct {
	ID                     string
	Name                   string // Имя пользователя
	Email                  string // Почта пользователя
	Role                   string // Роль пользователя
	CompletedVariantsCount int    // Количество сделанных вариантов
}
