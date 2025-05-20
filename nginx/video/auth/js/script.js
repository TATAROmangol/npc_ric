document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('errorMessage');
    const returnUrl = getReturnUrl();

    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const login = document.getElementById('login').value.trim();
            const password = document.getElementById('password').value;
            
            try {
                const response = await authenticateUser(login, password);
                
                if (response.token) {
                    // Сохраняем токен и перенаправляем
                    localStorage.setItem('authToken', response.token);
                    redirectAfterLogin(returnUrl);
                } else {
                    showError('Неверные учетные данные');
                }
            } catch (error) {
                console.error('Ошибка авторизации:', error);
                showError('Ошибка сервера. Попробуйте позже.');
            }
        });
    }

    // Если пользователь уже авторизован - редирект
    if (localStorage.getItem('authToken') && window.location.pathname === '/auth/') {
        redirectAfterLogin(returnUrl);
    }
});

// Функция аутентификации
async function authenticateUser(login, password) {
    const response = await fetch('/auth/api/login', {
        method: 'POST',
        credentials: 'include',  // Обязательно!
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login, password })
    });

    if (!response.ok) {
        throw new Error(await response.text());
    }

    // Проверяем cookie через 100ms
    return new Promise((resolve) => {
        setTimeout(() => {
            const token = document.cookie
                .split('; ')
                .find(row => row.startsWith('admin_token='))
                ?.split('=')[1];
            
            if (!token) throw new Error("Cookie not set");
            resolve(token);
        }, 300);
    });
}

// Вспомогательная функция для получения cookie
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

// Получение URL для возврата после авторизации
function getReturnUrl() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('returnUrl') || '/admin/';
}

// Редирект после успешной авторизации
function redirectAfterLogin(returnUrl) {
    window.location.href = returnUrl;
}

// Показ ошибки авторизации
function showError(message) {
    const errorElement = document.getElementById('errorMessage');
    if (errorElement) {
        errorElement.textContent = message;
        errorElement.style.display = 'block';
    }
}

// Выход из системы
function logout() {
    localStorage.removeItem('authToken');
    window.location.href = '/auth/';
}