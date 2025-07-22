// Элементы меню
const settingsBtn = document.getElementById('settingsBtn');
const closeMenuBtn = document.getElementById('closeMenu'); // Переименовано
const settingsMenu = document.getElementById('settingsMenu');
const overlay = document.getElementById('overlay');
const themeToggle = document.getElementById('themeToggle');

// Управление меню
settingsBtn.addEventListener('click', () => {
    settingsMenu.classList.add('open');
    overlay.classList.add('visible');
    document.body.classList.add('menu-open');
});

// Используем другое имя для функции закрытия
closeMenuBtn.addEventListener('click', closeSettingsMenu);
overlay.addEventListener('click', closeSettingsMenu);

function closeSettingsMenu() {
    settingsMenu.classList.remove('open');
    overlay.classList.remove('visible');
    document.body.classList.remove('menu-open');
}

// Закрытие по Esc
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') closeSettingsMenu();
});

// Управление темой
function setupTheme() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.body.classList.add('dark-theme');
        themeToggle.checked = true;
    }

    themeToggle.addEventListener('change', () => {
        document.body.classList.toggle('dark-theme', themeToggle.checked);
        localStorage.setItem('theme', themeToggle.checked ? 'dark' : 'light');
    });
}

// Проверка авторизации
async function checkAuth() {
    try {
        const response = await fetch('/api/auth/check', {
            method: 'GET',
            credentials: 'include'
        });
        
        if (!response.ok) return false;
        
        const result = await response.json();
        return result.authenticated === true;
    } catch (error) {
        console.error('Auth check failed:', error);
        return false;
    }
}

// Показ сообщения об авторизации
function showAuthMessage() {
    const authMessage = document.createElement('div');
    authMessage.className = 'auth-message';
    authMessage.textContent = 'Вы авторизованы';
    
    document.body.appendChild(authMessage);
    
    setTimeout(() => {
        authMessage.style.opacity = '0';
        setTimeout(() => authMessage.remove(), 500);
    }, 3000);
}

// Основная инициализация
async function init() {
    setupTheme();
    
    // Проверяем авторизацию только при загрузке страницы
    const isAuthenticated = await checkAuth();
    if (isAuthenticated) {
        showAuthMessage();
    }
}

// Запуск при полной загрузке DOM
document.addEventListener('DOMContentLoaded', init);