// Элементы меню
const settingsBtn = document.getElementById('settingsBtn');
const closeMenu = document.getElementById('closeMenu');
const settingsMenu = document.getElementById('settingsMenu');
const overlay = document.getElementById('overlay');
const themeToggle = document.getElementById('themeToggle');

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


window.addEventListener('DOMContentLoaded', loadTheme);

// Закрытие меню при нажатии Esc
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        settingsMenu.classList.remove('open');
        overlay.classList.remove('visible');
        document.body.classList.remove('menu-open');
    }
});