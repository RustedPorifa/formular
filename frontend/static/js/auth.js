document.addEventListener('DOMContentLoaded', function() {
    // Элементы интерфейса
    const tabs = document.querySelectorAll('.tab');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const passwordToggles = document.querySelectorAll('.toggle-password');
    const registerPassword = document.getElementById('registerPassword');
    const socialButtons = document.querySelectorAll('.social-btn');

    // Переключение между вкладками
    tabs.forEach(tab => {
        tab.addEventListener('click', function() {
            const targetTab = this.dataset.tab;
            
            // Обновляем активные вкладки
            tabs.forEach(t => t.classList.remove('active'));
            this.classList.add('active');
            
            // Показываем соответствующую форму
            document.querySelectorAll('.auth-form').forEach(form => {
                form.classList.remove('active');
            });
            document.getElementById(`${targetTab}Form`).classList.add('active');
        });
    });

    // Показать/скрыть пароль
    passwordToggles.forEach(toggle => {
        toggle.addEventListener('click', function() {
            const input = this.previousElementSibling;
            const type = input.type === 'password' ? 'text' : 'password';
            input.type = type;
        });
    });

    // Индикатор сложности пароля
    if (registerPassword) {
        registerPassword.addEventListener('input', function() {
            const password = this.value;
            const strengthBars = document.querySelectorAll('.password-strength .strength-bar');
            const hint = document.querySelector('.password-hint');
            
            // Сбросить стили
            strengthBars.forEach(bar => {
                bar.style.backgroundColor = '#e0e0e0';
                bar.style.flex = '1';
            });
            hint.style.display = 'none';

            if (password.length === 0) return;

            let strength = 0;
            const hasLetters = /[a-zA-Zа-яА-Я]/.test(password);
            const hasNumbers = /\d/.test(password);
            const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(password);

            if (password.length >= 8) strength++;
            if (hasLetters && hasNumbers) strength++;
            if (hasSpecial) strength++;

            // Обновить индикатор
            if (strength > 0) {
                for (let i = 0; i < strength; i++) {
                    strengthBars[i].style.backgroundColor = 
                        strength === 1 ? '#ff4d4d' : 
                        strength === 2 ? '#ffcc00' : '#4CAF50';
                }
            }
        });
    }

    // Обработка входа
    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const email = document.getElementById('loginEmail').value;
        const password = document.getElementById('loginPassword').value;
        const remember = this.querySelector('input[type="checkbox"]').checked;
        
        const button = this.querySelector('.auth-btn');
        button.disabled = true;
        button.textContent = 'Вход...';

        try {
            console.log(getCookie('XSRF-TOKEN'));
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'X-XSRF-TOKEN': getCookie('XSRF-TOKEN')},
                body: JSON.stringify({ email, password }),
                credentials: 'include',
            });

            const data = await response.json();
            
            if (response.ok) {
                // Перенаправление после успешного входа
                window.location.href = '/';
            } else {
                showError(this, data.error || 'Ошибка входа');
            }
        } catch (error) {
            showError(this, 'Сетевая ошибка');
        } finally {
            button.disabled = false;
            button.textContent = 'Войти';
        }
    });

    // Обработка регистрации
    registerForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const name = document.getElementById('registerName').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;
        const confirmPassword = document.getElementById('registerConfirm').value;
        const agreed = this.querySelector('input[type="checkbox"]').checked;
        
        // Валидация
        if (password !== confirmPassword) {
            showError(this, 'Пароли не совпадают');
            return;
        }
        
        if (!agreed) {
            showError(this, 'Необходимо согласие с условиями');
            return;
        }

        const button = this.querySelector('.auth-btn');
        button.disabled = true;
        button.textContent = 'Регистрация...';

        try {
            console.log(getCookie('XSRF-TOKEN'))
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'X-XSRF-TOKEN': getCookie('XSRF-TOKEN')},
                body: JSON.stringify({ name, email, password })
            });

            if (response.ok) {
                const data = await response.json();
                if (data.redirect) {
                    window.location.href = data.redirect;
                }
            } else {
                const data = await response.json();
                showError(this, data.error || 'Ошибка регистрации');
            }
        } catch (error) {
            showError(this, 'Сетевая ошибка');
        } finally {
            button.disabled = false;
            button.textContent = 'Зарегистрироваться';
        }
    });

    socialButtons.forEach(button => {
        button.addEventListener('click', function() {
            const provider = this.classList.contains('vk') ? 'VK' : 
            this.classList.contains('google') ? 'Google' : 'Yandex';
            alert(`Вход через ${provider} в разработке`);
        });
    });

    // Вспомогательные функции
    function showError(form, message) {
        let errorContainer = form.querySelector('.error-message');
        if (!errorContainer) {
            errorContainer = document.createElement('div');
            errorContainer.className = 'error-message';
            form.insertBefore(errorContainer, form.lastElementChild);
        }
        errorContainer.textContent = message;
        errorContainer.style.color = '#ff4d4d';
        errorContainer.style.marginTop = '10px';
        errorContainer.style.textAlign = 'center';
    }

    function showSuccess(message) {
        const successMessage = document.createElement('div');
        successMessage.className = 'success-message';
        successMessage.textContent = message;
        successMessage.style.color = '#4CAF50';
        successMessage.style.marginTop = '10px';
        successMessage.style.textAlign = 'center';
        successMessage.style.padding = '10px';
        
        const header = document.querySelector('.auth-header');
        header.appendChild(successMessage);
        
        setTimeout(() => {
            successMessage.remove();
        }, 5000);
    }
});

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}