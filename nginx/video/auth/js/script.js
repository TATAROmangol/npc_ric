document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('errorMessage');
    const successMessage = document.getElementById('successMessage');

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const login = document.getElementById('login').value;
        const password = document.getElementById('password').value;

        errorMessage.style.display = 'none';
        successMessage.style.display = 'none';

        try {
            const response = await fetch('/auth/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ login, password }),
                credentials: 'include' 
            });

            if (response.ok) {
                successMessage.textContent = 'Вход выполнен успешно! Перенаправление...';
                successMessage.style.display = 'block';
                setTimeout(() => window.location.href = '/admin/', 1500);
            } else {
                errorMessage.textContent = 'Неверный логин или пароль';
                errorMessage.style.display = 'block';
            }
        } catch (error) {
            errorMessage.textContent = 'Ошибка сети';
            errorMessage.style.display = 'block';
        }
    });

    async function checkAuthStatus() {
        try {
            const response = await fetch('/auth/api/login', {
                method: 'GET',
                credentials: 'include'
            });

            if (response.ok) {
                window.location.href = '/admin/';
            }
        } catch (error) {
            console.error('Ошибка проверки авторизации:', error);
        }
    }

    checkAuthStatus();
});