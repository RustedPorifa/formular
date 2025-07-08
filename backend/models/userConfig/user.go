package user

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"mail"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Completed []int  // Список ID выполненных вариантов
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
