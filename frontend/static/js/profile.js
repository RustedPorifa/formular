async function fetchUserProfile(isRetry = false) {
    try {
        const token = localStorage.getItem('jwtToken');
        if (!token) {
            throw new Error('Токен не найден');
        }

        const response = await fetch('http://127.0.0.1:8080/getuserinfo', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                //'Content-Type': 'application/json',
            },
            body: JSON.stringify({ token })
        });

        if (!response.ok) {
            const errorText = await response.text();
            // Добавляем статус в ошибку для последующей проверки
            const error = new Error(`HTTP ${response.status}: ${errorText}`);
            error.status = response.status;
            throw error;
        }

        const userData = await response.json();
        
        localStorage.setItem('userProfile', JSON.stringify({
            id: userData.id,
            name: userData.name,
            email: userData.email,
            role: userData.role,
            completed: userData.completed
        }));
        
        return userData;
    } catch (error) {
        // Проверяем 401 ошибку и отсутствие предыдущей попытки
        if (error.status === 401 && !isRetry) {
            try {
                console.log('Попытка обновления токена...');
                await refreshToken(); // Обновляем токен
                return fetchUserProfile(true); // Повторяем запрос с флагом isRetry
            } catch (refreshError) {
                console.error('Ошибка обновления токена:', refreshError);
                throw refreshError;
            }
        }
        
        // Все другие ошибки или повторная 401 ошибка
        console.error('Ошибка:', error);
        throw error;
    }
}

async function refreshToken() {
    try {
        const response = await fetch('http://127.0.0.1:8080/refreshtokens', {
            method: 'POST',
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error(`Ошибка ${response.status} при обновлении токена`);
        }

        const data = await response.json();
        // Исправлено: единый ключ для токена
        localStorage.setItem('jwtToken', data.access_token);
        return data.access_token;
    } catch (error) {
        console.error("Refresh error:", error);
        window.location.href = "/";
        throw error;
    }
}

// Пример использования:
async function loadUserProfile() {
    try {
        const profile = await fetchUserProfile();
        console.log('Данные пользователя:', profile);
        // Обновляем UI
        document.getElementById('user-name').textContent = profile.name;
        document.getElementById('user-email').textContent = profile.email;
        // ... и т.д.
    } catch (error) {
        console.error('Не удалось загрузить профиль:', error);
        alert(`Ошибка: ${error.message}`);
    }
}



// Вызываем при загрузке страницы профиля
document.addEventListener('DOMContentLoaded', loadUserProfile);