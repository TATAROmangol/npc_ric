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

            if (response.status === 200) {
                setTimeout(() => window.location.href = '/admin/', 1000);
            } 
            
        } catch (error) {
            errorMessage.textContent = 'Ошибка сети';
            errorMessage.style.display = 'block';
        }
    });
});