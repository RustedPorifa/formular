// Элементы меню
const settingsBtn = document.getElementById('settingsBtn');
const closeMenu = document.getElementById('closeMenu');
const settingsMenu = document.getElementById('settingsMenu');
const overlay = document.getElementById('overlay');
const themeToggle = document.getElementById('themeToggle');

async function checkAuthAndUpdateProfile() {
    try {
        const response = await fetch('/api/verify', {
            method: 'GET',
            credentials: 'include' // Важно для отправки cookies
        });

        if (response.ok) {
            const data = await response.json();
            if (data.verify === "true") {
                updateProfileButton();
            }
        }
    } catch (error) {
        console.error('Ошибка при проверке аутентификации:', error);
    }
}

// Обновление кнопки профиля
function updateProfileButton() {
    const authBtn = document.getElementById('auth-profile-btn');
    if (!authBtn) return;


    authBtn.onclick = function() {
        window.location = '/user/profile';
    };
    authBtn.textContent = 'Профиль';
}

// Открытие меню
settingsBtn.addEventListener('click', () => {
    settingsMenu.classList.add('open');
    overlay.classList.add('visible');
    document.body.classList.add('menu-open');
});

// Закрытие меню
closeMenu.addEventListener('click', () => {
    settingsMenu.classList.remove('open');
    overlay.classList.remove('visible');
    document.body.classList.remove('menu-open');
});

// Закрытие по клику на overlay
overlay.addEventListener('click', () => {
    settingsMenu.classList.remove('open');
    overlay.classList.remove('visible');
    document.body.classList.remove('menu-open');
});

// Сохранение и загрузка темы
const loadTheme = () => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.body.classList.add('dark-theme');
        themeToggle.checked = true;
    }
};


themeToggle.addEventListener('change', () => {
    if (themeToggle.checked) {
        document.body.classList.add('dark-theme');
        localStorage.setItem('theme', 'dark');
    } else {
        document.body.classList.remove('dark-theme');
        localStorage.setItem('theme', 'light');
    }
});


window.addEventListener('DOMContentLoaded', () => {
    loadTheme();
    checkAuthAndUpdateProfile(); // Добавлен вызов проверки аутентификации
});

// Закрытие меню при нажатии Esc
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        settingsMenu.classList.remove('open');
        overlay.classList.remove('visible');
        document.body.classList.remove('menu-open');
    }
});