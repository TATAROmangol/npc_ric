class ApiService {
    async fetchWithAuth(url, options = {}) {
        const token = localStorage.getItem('authToken');
        const headers = {
            'Content-Type': 'application/json',
            ...(token && { 'Authorization': `Bearer ${token}` }),
            ...options.headers
        };

        const response = await fetch(url, {
            ...options,
            headers
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return response.json();
    }

    async getMentors() {
        return this.fetchWithAuth('/admin/get/mentors');
    }

    async addMentor(name) {
        return this.fetchWithAuth('/admin/post/mentor', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name })
        });
    }

    async deleteMentor(id) {
        return this.fetchWithAuth('/admin/delete/mentor', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id })
        });
    }

    async updateMentor(id, newName) {
    return this.fetchWithAuth('/admin/put/mentor', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id, name: newName })
    });
}
}

const apiService = new ApiService();
let mentors = [];
let selectedMentor = null;

// Получаем элементы DOM
const supervisorsList = document.getElementById('supervisorsList');
const openSupervisorModalBtn = document.getElementById('openSupervisorModalBtn');
const addSupervisorModal = document.getElementById('addSupervisorModal');
const supervisorNameInput = document.getElementById('supervisorName');
const submitSupervisorBtn = document.getElementById('submitSupervisorBtn');
const cancelSupervisorBtn = document.getElementById('cancelSupervisorBtn');
const closeSupervisorBtn = document.querySelector('.closeSupervisorBtn');
const deleteSupervisorBtn = document.getElementById('deleteSupervisorBtn');

// Инициализация
async function init() {
    try {
        await loadMentors();
        
        // Обработчики событий
        openSupervisorModalBtn.addEventListener('click', openAddSupervisorModal);
        submitSupervisorBtn.addEventListener('click', addSupervisor);
        cancelSupervisorBtn.addEventListener('click', closeAddSupervisorModal);
        closeSupervisorBtn.addEventListener('click', closeAddSupervisorModal);
        deleteSupervisorBtn.addEventListener('click', deleteSupervisor);
        
        // Закрытие при клике вне окна
        window.addEventListener('click', (e) => {
            if (e.target === addSupervisorModal) closeAddSupervisorModal();
        });
    } catch (error) {
        console.error('Ошибка инициализации:', error);
        alert('Не удалось загрузить данные руководителей');
    }
}

// Загрузка руководителей с сервера
async function loadMentors() {
    try {
        mentors = await apiService.getMentors();
        renderSupervisorsList();
    } catch (error) {
        console.error('Ошибка загрузки руководителей:', error);
        throw error;
    }
}

// Рендер списка
function renderSupervisorsList() {
    supervisorsList.innerHTML = '';
    mentors.forEach(mentor => {
        const item = document.createElement('div');
        item.className = 'list-item';
        item.textContent = mentor.name;
        item.dataset.id = mentor.id;
        item.addEventListener('click', () => selectSupervisor(mentor));
        supervisorsList.appendChild(item);
    });
}

// Выбор руководителя
function selectSupervisor(mentor) {
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
function openAddSupervisorModal() {
    console.log('Функция openAddSupervisorModal вызвана'); // Добавьте это
    supervisorNameInput.value = '';
    addSupervisorModal.style.display = 'block';
}

function closeAddSupervisorModal() {
    addSupervisorModal.style.display = 'none';
}

async function addSupervisor() {
    const name = supervisorNameInput.value.trim();
    if (!name) {
        alert('Введите ФИО руководителя');
        return;
    }

    try {
        const newMentor = await apiService.addMentor(name);
        mentors.push(newMentor);
        renderSupervisorsList();
        closeAddSupervisorModal();
        alert('Руководитель успешно добавлен!');
    } catch (error) {
        console.error('Ошибка добавления руководителя:', error);
        alert('Не удалось добавить руководителя');
    }
}

// Удаление руководителя
async function deleteSupervisor() {
    if (!selectedMentor) {
        alert('Выберите руководителя для удаления');
        return;
    }
    
    if (confirm(`Удалить руководителя "${selectedMentor.name}"?`)) {
        try {
            await apiService.deleteMentor(selectedMentor.id);
            mentors = mentors.filter(m => m.id !== selectedMentor.id);
            selectedMentor = null;
            renderSupervisorsList();
            alert('Руководитель успешно удален');
        } catch (error) {
            console.error('Ошибка удаления руководителя:', error);
            alert('Не удалось удалить руководителя');
        }
    }
}

// Запускаем приложение
document.addEventListener('DOMContentLoaded', init);