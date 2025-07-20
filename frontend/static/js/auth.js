// Переключение между вкладками входа и регистрации
document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', () => {
        // Удаляем активный класс у всех вкладок
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        
        // Добавляем активный класс текущей вкладке
        tab.classList.add('active');
        
        // Скрываем все формы
        document.querySelectorAll('.auth-form').forEach(form => {
            form.classList.remove('active');
        });
        
        // Показываем соответствующую форму
        const tabName = tab.getAttribute('data-tab');
        document.getElementById(`${tabName}Form`).classList.add('active');
    });
});

// Переключение видимости пароля
document.querySelectorAll('.toggle-password').forEach(button => {
    button.addEventListener('click', function() {
        const input = this.parentElement.querySelector('input');
        const type = input.getAttribute('type') === 'password' ? 'text' : 'password';
        input.setAttribute('type', type);
        
        // Меняем иконку
        this.textContent = type === 'password' ? '👁️' : '🔒';
    });
});

// Проверка сложности пароля при регистрации
const passwordInput = document.getElementById('registerPassword');
if (passwordInput) {
    passwordInput.addEventListener('input', function() {
        const password = this.value;
        const strengthBars = document.querySelectorAll('.strength-bar');
        
        // Сбрасываем все бары
        strengthBars.forEach(bar => {
            bar.style.background = '#e9ecef';
            if (document.body.classList.contains('dark-theme')) {
                bar.style.background = '#444';
            }
        });
        
        // Проверяем сложность пароля
        if (password.length > 0) {
            strengthBars[0].style.background = 'var(--danger)';
            
            if (password.length >= 6) {
                strengthBars[1].style.background = 'var(--warning)';
            }
            
            if (password.length >= 8 && /[A-Z]/.test(password) && /[0-9]/.test(password)) {
                strengthBars[2].style.background = 'var(--success)';
            }
        }
    });
}

// Обработка формы входа
document.getElementById('loginForm')?.addEventListener('submit', async function(e) {
    e.preventDefault();
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;
    
    await login(email, password);
});

// Обработка формы регистрации
document.getElementById('registerForm')?.addEventListener('submit', async function(e) {
    e.preventDefault();
    const name = document.getElementById('registerName').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;
    const confirmPassword = document.getElementById('registerConfirm').value;
    
    // Проверка совпадения паролей
    if (password !== confirmPassword) {
        alert('Пароли не совпадают!');
        return;
    }
    
    // Проверка сложности пароля
    if (password.length < 8) {
        alert('Пароль должен содержать не менее 8 символов');
        return;
    }
    
    await register(name, email, password);
    // Переключаем на вкладку входа после успешной регистрации
    document.querySelector('[data-tab="login"]').click();
});

// Функция входа
async function login(email, password) {
    try {
        const response = await fetch('http://127.0.0.1:8080/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Ошибка сервера');
        }

        const data = await response.json();
        localStorage.setItem('jwtToken', data.token);
        alert('Успешный вход!');
        // Дальнейшие действия после входа
        // window.location.href = "/dashboard.html";
    } catch (error) {
        console.error("Login error:", error);
        alert(`Ошибка входа: ${error.message}`);
    }
}

// Функция регистрации
async function register(name, email, password) {
    try {
        const response = await fetch('http://127.0.0.1:8080/register', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ name, email, password })
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Ошибка сервера');
        }

        alert('Регистрация успешна! Теперь войдите');
    } catch (error) {
        console.error("Registration error:", error);
        alert(`Ошибка регистрации: ${error.message}`);
    }
}

// Инициализация темы
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.body.classList.add('dark-theme');
    }
    
    // Применяем тему к элементам, которые зависят от темы
    if (document.body.classList.contains('dark-theme')) {
        const strengthBars = document.querySelectorAll('.strength-bar');
        if (strengthBars.length > 0 && strengthBars[0].style.background === '') {
            strengthBars.forEach(bar => {
                bar.style.background = '#444';
            });
        }
    }
});