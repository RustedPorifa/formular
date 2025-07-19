// script.js
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
    const mobileIcon = document.getElementById('mobile-theme-icon');
    if (theme === 'dark') {
        icon.innerHTML = `<path d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />`;
        mobileIcon.innerHTML = `<path d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />`;
    } else {
        icon.innerHTML = `<path d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />`;
        mobileIcon.innerHTML = `<path d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />`;
    }
}

// Мобильное меню
function toggleMobileMenu() {
    const menu = document.getElementById('mobile-menu');
    menu.classList.toggle('hidden');
    const button = document.getElementById('mobile-menu-toggle');
    button.innerHTML = menu.classList.contains('hidden') ? 
        `<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M4 6h16M4 12h16M4 18h16" />
        </svg>` :
        `<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M6 18L18 6M6 6l12 12" />
        </svg>`;
}

// Mock data для вариантов ЕГЭ
const examVariants = [
    { id: 1, title: "Вариант 1", difficulty: "easy"},
    { id: 2, title: "Вариант 2", difficulty: "medium"},
    { id: 3, title: "Вариант 3", difficulty: "hard"},
    { id: 4, title: "Вариант 4", difficulty: "medium"},
    { id: 5, title: "Вариант 5", difficulty: "easy"},
    { id: 6, title: "Вариант 6", difficulty: "medium"},
    { id: 7, title: "Вариант 7", difficulty: "hard"},
    { id: 8, title: "Вариант 8", difficulty: "easy"},
    { id: 9, title: "Вариант 9", difficulty: "medium"},
];

// Создание карточки варианта
function createExamCard(variant) {
    const card = document.createElement('div');
    card.className = 'card animate-fade-in';
    card.style.opacity = '0';
    card.style.transform = 'translateY(20px)';
    
    card.innerHTML = `
        <div class="difficulty-bar ${variant.difficulty}"></div>
        <div class="card-body">
            <div class="flex justify-between items-start">
                <h3 class="card-title">${variant.title}</h3>
                <span class="difficulty-badge ${variant.difficulty}">${getDifficultyLabel(variant.difficulty)}</span>
            </div>
            <div class="space-y-2 mt-3">
                <div class="info-item">
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                    </svg>
                    Вопросов: ${variant.questions}
                </div>
                <div class="info-item">
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Длительность: ${variant.duration}
                </div>
            </div>
            <div class="card-footer">
                <button class="start-button">Начать практику</button>
            </div>
        </div>
    `;
    
    return card;
}

// Получение текстовой метки уровня сложности
function getDifficultyLabel(difficulty) {
    switch(difficulty) {
        case 'easy': return 'Легкий';
        case 'medium': return 'Средний';
        case 'hard': return 'Сложный';
        default: return 'Неизвестно';
    }
}

// Отрисовка карточек
function renderExamCards(variants) {
    const container = document.getElementById('exam-variants');
    container.innerHTML = '';
    
    if (variants.length === 0) {
        container.innerHTML = `
            <div class="empty-state animate-fade-in" style="opacity: 0; transform: translateY(20px)">
                <svg viewBox="0 0 24 24" width="64" height="64" fill="none" stroke="#93C5FD" stroke-width="2" class="mx-auto mb-4">
                    <path d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <h3 class="text-lg font-medium mb-1">Варианты не найдены</h3>
                <p class="mb-4">Попробуйте выбрать другой уровень сложности.</p>
                <button id="show-all" class="btn btn-primary">Показать все уровни</button>
            </div>
        `;
        
        document.getElementById('show-all').addEventListener('click', () => {
            document.querySelectorAll('.filter-button').forEach(btn => btn.classList.remove('active'));
            document.querySelector('.filter-button.all').classList.add('active');
            renderExamCards(examVariants);
        });
        
        setTimeout(() => {
            container.querySelector('.empty-state').style.opacity = 1;
            container.querySelector('.empty-state').style.transform = 'translateY(0)';
        }, 100);
        
        return;
    }
    
    variants.forEach((variant, index) => {
        const card = createExamCard(variant);
        container.appendChild(card);
        
        // Анимация появления с задержкой
        setTimeout(() => {
            card.style.opacity = 1;
            card.style.transform = 'translateY(0)';
        }, index * 100);
    });
}

// Фильтрация вариантов
function setupFilters() {
    const filters = {
        all: () => examVariants,
        easy: () => examVariants.filter(v => v.difficulty === 'easy'),
        medium: () => examVariants.filter(v => v.difficulty === 'medium'),
        hard: () => examVariants.filter(v => v.difficulty === 'hard')
    };
    
    Object.entries(filters).forEach(([key, filterFn]) => {
        document.querySelector(`.filter-button.${key}`).addEventListener('click', () => {
            document.querySelectorAll('.filter-button').forEach(btn => btn.classList.remove('active'));
            document.querySelector(`.filter-button.${key}`).classList.add('active');
            renderExamCards(filterFn());
        });
    });
}

// Инициализация
async function init() {
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
    document.getElementById('mobile-theme-toggle').addEventListener('click', toggleTheme);
    document.getElementById('mobile-menu-toggle').addEventListener('click', toggleMobileMenu);
    document.getElementById("login-profile-button").addEventListener('click', GetProfile)
    // Инициализация фильтров и карточек
    setupFilters();
    renderExamCards(examVariants);
    await changeProfileButton();
    // Активная кнопка по умолчанию
    document.querySelector('.filter-button.all').classList.add('active');
}

async function changeProfileButton() {
    if (isLocalStorageAvailable()) {
        let jwt = localStorage.getItem("jwtToken");
        console.log("JWT:", jwt);
        if (jwt !== null) {
            let button = document.getElementById("login-profile-button")
            let attributes = button.attributes;
            for (let i = attributes.length - 1; i >= 0; i--) {
                button.removeAttribute(attributes[i].name);
            }
            button.textContent = "Профиль"
        }
    } else {
        console.log("Local storage was blocked")
    }
}

async function GetProfile() {
    // Проверяем доступность localStorage
    if (!isLocalStorageAvailable()) {
        redirectToLogin();
        return;
    }

    // Получаем JWT из localStorage
    const jwt = localStorage.getItem("jwtToken");
    if (!jwt) {
        redirectToLogin();
        return;
    }

    // 1. Сначала получаем данные пользователя
    try {
        const response = await fetch("/getuserinfo", {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${jwt}`
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                redirectToLogin();
            } else {
                console.error("Server error:", response.status);
            }
            return;
        }

        const profileData = await response.json();
        
        // 2. Сохраняем данные в sessionStorage
        sessionStorage.setItem("profileData", JSON.stringify(profileData));
        
        // 3. Переходим на страницу профиля
        window.location.href = '/profile';
        
    } catch (error) {
        console.error("Error fetching profile:", error);
        redirectToLogin();
    }
}

function redirectToLogin() {
    window.location.href = '/loginform';
}

// Запуск инициализации после загрузки DOM
document.addEventListener('DOMContentLoaded', init);