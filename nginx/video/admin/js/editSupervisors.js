class ApiService {
    constructor() {
        this.adminBaseUrl = '/admin/api';
        this.formsBaseUrl = '/forms/api';
    }
    
    async fetchWithAuth(url, options = {}, expectJson = true) {
        console.log('Отправка запроса:', { url, options });
        const token = localStorage.getItem('authToken');
        if (token) {
            options.headers = {
                ...(options.headers || {}),
                'Authorization': `Bearer ${token}`
            };
        }
        
        const response = await fetch(url, options);
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Ошибка ответа:', errorText);
            throw new Error(`Ошибка сервера: ${response.status} - ${errorText}`);
        }
        
        // Если не ожидаем JSON, просто возвращаем response
        return expectJson ? response.json() : response;
    }

    // Mentor methods
    async getMentors() {
        return this.fetchWithAuth(`${this.adminBaseUrl}/get/mentors`);
    }

    async addMentor(name) {
        if (!name || typeof name !== 'string') {
            throw new Error('Имя руководителя обязательно и должно быть строкой');
        }

        return this.fetchWithAuth(`${this.adminBaseUrl}/post/mentor`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                name: name
            })
        });
    }

    async updateMentor(data) {
        return this.fetchWithAuth(`${this.adminBaseUrl}/put/mentor`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
    }

    async deleteMentor(id) {
        if (!id) {
            throw new Error('ID руководителя обязательно');
        }

        const response = await this.fetchWithAuth(`${this.adminBaseUrl}/delete/mentor`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id })
        }, false); // Добавляем флаг, что не ожидаем JSON в ответе
    }
}

const apiService = new ApiService();
let mentors = [];
let selectedMentor = null;
let mentorsList;
let openMentorModalBtn, addMentorModal, mentorNameInput;
let submitMentorBtn, cancelMentorBtn, closeMentorBtn, deleteMentorBtn;

// Инициализация
async function init() {
    // Инициализация DOM-элементов
    mentorsList = document.getElementById('mentorsList');
    openMentorModalBtn = document.getElementById('openMentorModalBtn');
    addMentorModal = document.getElementById('addMentorModal');
    mentorNameInput = document.getElementById('mentorName');
    submitMentorBtn = document.getElementById('submitMentorBtn');
    cancelMentorBtn = document.getElementById('cancelMentorBtn');
    closeMentorBtn = document.querySelector('.closeMentorBtn');
    deleteMentorBtn = document.getElementById('deleteMentorBtn');
    if (!mentorsList) {
        console.error('Элемент mentorsList не найден в DOM');
        return;
    }
    
    try {
        await loadMentors();
        
        // Обработчики событий
        openMentorModalBtn.addEventListener('click', openAddMentorModal);
        submitMentorBtn.addEventListener('click', addMentor);
        cancelMentorBtn.addEventListener('click', closeAddMentorModal);
        closeMentorBtn.addEventListener('click', closeAddMentorModal);
        deleteMentorBtn.addEventListener('click', deleteMentor);
        
        window.addEventListener('click', (e) => {
            if (e.target === addMentorModal) closeAddMentorModal();
        });
    } catch (error) {
        console.error('Ошибка инициализации:', error);
        showCustomAlert("Не удалось загрузить данные руководителей");
    }
}

// Загрузка руководителей с сервера
async function loadMentors() {
    try {
        const response = await apiService.getMentors();
        mentors = Array.isArray(response) ? response : [];
        renderMentorsList();
    } catch (error) {
        console.error('Ошибка загрузки руководителей:', error);
        mentors = [];
        renderMentorsList();
        throw error;
    }
}

// Рендер списка
function renderMentorsList() {
    if (!mentorsList) {
        console.error('Элемент mentorsList не инициализирован');
        return;
    }
    
    mentorsList.innerHTML = '';
    
    if (!mentors || !Array.isArray(mentors)) {
        mentorsList.innerHTML = '<div class="empty">Нет данных о руководителях</div>';
        return;
    }
    
    mentors.forEach(mentor => {
        const item = document.createElement('div');
        item.className = 'list-item';
        item.textContent = mentor.name;
        item.dataset.id = mentor.id;
        item.addEventListener('click', () => selectMentor(mentor));
        mentorsList.appendChild(item);
    });
}

// Выбор руководителя
function selectMentor(mentor) {
    document.querySelectorAll('.list-item').forEach(item => {
        item.classList.remove('selected');
    });
    
    const selectedItem = [...document.querySelectorAll('.list-item')]
        .find(item => item.dataset.id === mentor.id.toString());
    
    if (selectedItem) {
        selectedItem.classList.add('selected');
        selectedMentor = mentor;
    }
}

// Работа с модальным окном
function openAddMentorModal() {
    console.log('Функция openAddMentorModal вызвана'); 
    mentorNameInput.value = '';
    addMentorModal.style.display = 'block';
}

function closeAddMentorModal() {
    addMentorModal.style.display = 'none';
}

async function addMentor() {
    const name = mentorNameInput.value.trim();
    
    if (!name) {
        showCustomAlert("Введите ФИО руководителя");
        return;
    }

    try {
        const newMentor = await apiService.addMentor(name);
        mentors.push(newMentor);
        renderMentorsList();
        closeAddMentorModal();
        showCustomAlert("Руководитель успешно добавлен");
    } catch (error) {
        console.error('Ошибка добавления руководителя:', error);
        showCustomAlert(`Ошибка: ${error.message}`);
    }
}

// Удаление руководителя
async function deleteMentor() {
    if (!selectedMentor) {
        showCustomAlert("Выберите руководителя для удаления");
        return;
    }
    
    if (confirm(`Удалить руководителя "${selectedMentor.name}"?`)) {
        try {
            await apiService.deleteMentor(selectedMentor.id);
            mentors = mentors.filter(m => m.id !== selectedMentor.id);
            selectedMentor = null;
            renderMentorsList();
            showCustomAlert("Руководитель успешно удалён");
        } catch (error) {
            console.error('Ошибка удаления руководителя:', error);
            showCustomAlert(`Ошибка удаления: ${error.message}`);
        }
    }
}

function showCustomAlert(message) {
    document.getElementById('customAlertMessage').textContent = message;
    document.getElementById('customAlert').style.display = 'block';
    document.getElementById('customAlertOverlay').style.display = 'block';
}

function hideCustomAlert() {
    document.getElementById('customAlert').style.display = 'none';
    document.getElementById('customAlertOverlay').style.display = 'none';
}

// Запускаем приложение
document.addEventListener('DOMContentLoaded', init);