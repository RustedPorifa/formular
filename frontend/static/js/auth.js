// auth.js
// Проверка доступности localStorage
function isLocalStorageAvailable() {
    try {
        const testKey = '__test__storage__';
        localStorage.setItem(testKey, testKey);
        localStorage.removeItem(testKey);
        return true;
    } catch (e) {
        return false;
    }
}

// Переключение темы
function toggleTheme() {
    const isDark = document.body.classList.contains('dark');
    if (isDark) {
        document.body.classList.remove('dark');
        document.body.classList.add('light');
        updateThemeIcon('light');
    } else {
        document.body.classList.remove('light');
        document.body.classList.add('dark');
        updateThemeIcon('dark');
    }
    
    if (isLocalStorageAvailable()) {
        localStorage.setItem('theme', isDark ? 'light' : 'dark');
    }
}

// Обновление иконки темы
function updateThemeIcon(theme) {
    const icon = document.getElementById('theme-icon');
    if (theme === 'dark') {
        icon.innerHTML = `<path d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />`;
    } else {
        icon.innerHTML = `<path d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />`;
    }
}

// Переключение между формами
function setupFormSwitch() {
    const loginTab = document.getElementById('login-tab');
    const registerTab = document.getElementById('register-tab');
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    
    loginTab.addEventListener('click', () => {
        loginTab.classList.add('active');
        registerTab.classList.remove('active');
        loginForm.classList.add('active');
        registerForm.classList.remove('active');
    });
    
    registerTab.addEventListener('click', () => {
        registerTab.classList.add('active');
        loginTab.classList.remove('active');
        registerForm.classList.add('active');
        loginForm.classList.remove('active');
    });
}

async function login(params) {
    
}
async function register() {
    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({name, email, password})
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message);
        }

        alert('Регистрация успешна! Теперь войдите');
        
    } catch (error) {
        alert(`Ошибка: ${error.message}`);
    }
}

async function getProtectedData() {
    try {
        // Получаем токен из localStorage
        const token = localStorage.getItem('jwtToken');
        
        const response = await fetch('http://localhost:8080/protected', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Доступ запрещен');
        }

        const data = await response.json();
        console.log('Защищенные данные:', data);
        
    } catch (error) {
        alert(`Ошибка: ${error.message}`);
    }
}

// Валидация форм
function setupFormValidation() {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    
    loginForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;
        
        if (email && password) {
            // Здесь должна быть логика отправки данных
            alert('Вход выполнен!');
        } else {
            alert('Пожалуйста, заполните все поля');
        }
    });
    
    registerForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const name = document.getElementById('register-name').value;
        const email = document.getElementById('register-email').value;
        const password = document.getElementById('register-password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        
        if (name && email && password && confirmPassword) {
            if (password === confirmPassword) {
                // Здесь должна быть логика регистрации
                alert('Регистрация успешна!');
            } else {
                alert('Пароли не совпадают');
            }
        } else {
            alert('Пожалуйста, заполните все поля');
        }
    });
}

// Инициализация
function init() {
    // Проверка темы
    if (isLocalStorageAvailable()) {
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme) {
            document.body.className = savedTheme;
            updateThemeIcon(savedTheme);
        }
    }
    
    // Обработчики событий
    document.getElementById('theme-toggle').addEventListener('click', toggleTheme);
    setupFormSwitch();
    setupFormValidation();
}

// Запуск инициализации после загрузки DOM
document.addEventListener('DOMContentLoaded', init);
document.getElementById('login-btn').addEventListener('click', login);
document.getElementById('register-btn').addEventListener('click', register);
document.getElementById('data-btn').addEventListener('click', getProtectedData);