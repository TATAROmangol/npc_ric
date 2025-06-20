document.addEventListener('DOMContentLoaded', function() {
    // Получаем элементы DOM
    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('errorMessage');
    const successMessage = document.getElementById('successMessage');
    const loginInput = document.getElementById('login');
    const passwordInput = document.getElementById('password');
    const submitButton = document.getElementById('submitButton');

    // Обработчик отправки формы
    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        // Получаем значения полей
        const login = loginInput.value.trim();
        const password = passwordInput.value;

        // Сбрасываем сообщения
        errorMessage.style.display = 'none';
        successMessage.style.display = 'none';
        
        // Валидация полей
        if (!login || !password) {
            showAuthError('Заполните все обязательные поля');
            return;
        }

        try {
            // Отправляем запрос на сервер
            const response = await fetch('/auth/api/login', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify({ login, password }),
                credentials: 'include' // Для работы с куками
            });

            // Обрабатываем ответ
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || 'Ошибка авторизации');
            }

            // Успешная авторизация
            showAuthSuccess('Авторизация успешна! Перенаправление...');
            
            // Перенаправляем в админку
            setTimeout(() => {
                window.location.href = '/admin/';
            }, 1000);
            
        } catch (error) {
            console.error('Ошибка авторизации:', error);
            showAuthError(error.message || 'Ошибка сети при авторизации');
        } finally {
            // Разблокируем кнопку
            submitButton.disabled = false;
            submitButton.textContent = 'Войти';
        }
    });

    // Функция показа ошибки авторизации
    function showAuthError(message) {
        errorMessage.textContent = message;
        errorMessage.style.display = 'block';
        errorMessage.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }

    // Функция показа успешного сообщения
    function showAuthSuccess(message) {
        successMessage.textContent = message;
        successMessage.style.display = 'block';
        successMessage.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }

    // Очистка сообщений при изменении полей
    loginInput.addEventListener('input', () => {
        errorMessage.style.display = 'none';
    });
    
    passwordInput.addEventListener('input', () => {
        errorMessage.style.display = 'none';
    });
});