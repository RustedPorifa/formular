// Элементы меню
const settingsBtn = document.getElementById('settingsBtn');
const closeMenu = document.getElementById('closeMenu');
const settingsMenu = document.getElementById('settingsMenu');
const overlay = document.getElementById('overlay');
const themeToggle = document.getElementById('themeToggle');
const authProfileBtn = document.getElementById('auth-profile-btn');




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

// Переключение темы
themeToggle.addEventListener('change', () => {
    if (themeToggle.checked) {
        document.body.classList.add('dark-theme');
        localStorage.setItem('theme', 'dark');
    } else {
        document.body.classList.remove('dark-theme');
        localStorage.setItem('theme', 'light');
    }
});

function init() {
    loadTheme()

}

function changeProfileButton() {
    let refreshToken = getCookie("access_token");
    if (refreshToken !== "") {
        authProfileBtn.textContent = "Профиль"
        authProfileBtn.onclick="location.href='profile'"
    } else {
        authProfileBtn.onclick="location.href='loginform'"
    }
}

function getCookie(name) {
  const nameEQ = name + "=";
  const ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') c = c.substring(1, c.length); // Remove leading spaces
    if (c.indexOf(nameEQ) === 0) {
      return c.substring(nameEQ.length, c.length); // Return the cookie value
    }
  }
  return null; // Return null if the cookie is not found
}

// Инициализация темы при загрузке
window.addEventListener('DOMContentLoaded', init());

// Закрытие меню при нажатии Esc
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        settingsMenu.classList.remove('open');
        overlay.classList.remove('visible');
        document.body.classList.remove('menu-open');
    }
});