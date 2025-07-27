// Элементы DOM
const emailModal = document.getElementById('emailModal');
const passwordModal = document.getElementById('passwordModal');
const tariffsModal = document.getElementById('tariffsModal');
const changeEmailBtn = document.getElementById('changeEmailBtn');
const changePasswordBtn = document.getElementById('changePasswordBtn');
const viewTariffsBtn = document.getElementById('viewTariffsBtn');
const closeModalBtns = document.querySelectorAll('.close-modal');
const logoutBtn = document.getElementById('logoutBtn');
const settingsBtn = document.getElementById('settingsBtn');

// Открытие модальных окон
changeEmailBtn.addEventListener('click', () => {
    emailModal.classList.add('show');
    document.body.classList.add('modal-open');
});

changePasswordBtn.addEventListener('click', () => {
    passwordModal.classList.add('show');
    document.body.classList.add('modal-open');
});

viewTariffsBtn.addEventListener('click', () => {
    tariffsModal.classList.add('show');
    document.body.classList.add('modal-open');
});

// Закрытие модальных окон
closeModalBtns.forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('.modal').forEach(modal => {
            modal.classList.remove('show');
        });
        document.body.classList.remove('modal-open');
    });
});

// Закрытие по клику вне модального окна
window.addEventListener('click', (e) => {
    if (e.target.classList.contains('modal')) {
        e.target.classList.remove('show');
        document.body.classList.remove('modal-open');
    }
});

// Обработка форм
document.getElementById('emailForm')?.addEventListener('submit', function(e) {
    e.preventDefault();
    const newEmail = document.getElementById('newEmail').value;
    
    // Обновляем email в профиле
    document.getElementById('userEmail').textContent = newEmail;
    
    // Закрываем модальное окно
    emailModal.classList.remove('show');
    document.body.classList.remove('modal-open');
    
    // Показываем уведомление
    alert('Email успешно изменен!');
});

document.getElementById('passwordForm')?.addEventListener('submit', function(e) {
    e.preventDefault();
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    
    // Проверка совпадения паролей
    if (newPassword !== confirmPassword) {
        alert('Пароли не совпадают!');
        return;
    }
    
    // Проверка сложности пароля
    if (newPassword.length < 8) {
        alert('Пароль должен содержать не менее 8 символов');
        return;
    }
    
    // Закрываем модальное окно
    passwordModal.classList.remove('show');
    document.body.classList.remove('modal-open');
    
    // Показываем уведомление
    alert('Пароль успешно изменен!');
});

// Индикатор сложности пароля
const passwordInput = document.getElementById('newPassword');
if (passwordInput) {
    passwordInput.addEventListener('input', function() {
        const password = this.value;
        const strengthBars = document.querySelectorAll('#passwordModal .strength-bar');
        
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

// Выход из аккаунта
logoutBtn.addEventListener('click', () => {
    if (confirm('Вы уверены, что хотите выйти из аккаунта?')) {
        // Перенаправляем на главную страницу
        window.location.href = 'index.html';
    }
});

// Инициализация темы
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.body.classList.add('dark-theme');
    }
});