let supervisors = [
    "Иванов Иван Иванович",
    "Петров Петр Петрович"
];

let selectedSupervisor = null;

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
function init() {
    renderSupervisorsList();
    
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
}

// Рендер списка
function renderSupervisorsList() {
    supervisorsList.innerHTML = '';
    supervisors.forEach(supervisor => {
        const item = document.createElement('div');
        item.className = 'list-item';
        item.textContent = supervisor;
        item.addEventListener('click', () => selectSupervisor(supervisor));
        supervisorsList.appendChild(item);
    });
}

// Выбор руководителя
function selectSupervisor(supervisorName) {
    document.querySelectorAll('.list-item').forEach(item => {
        item.classList.remove('selected');
    });
    
    const selectedItem = [...document.querySelectorAll('.list-item')]
        .find(item => item.textContent === supervisorName);
    
    if (selectedItem) {
        selectedItem.classList.add('selected');
        selectedSupervisor = supervisorName;
    }
}

// Работа с модальным окном
function openAddSupervisorModal() {
    console.log('Открываем модальное окно');
    supervisorNameInput.value = '';
    addSupervisorModal.style.display = 'block';
}

function closeAddSupervisorModal() {
    addSupervisorModal.style.display = 'none';
}

function addSupervisor() {
    const name = supervisorNameInput.value.trim();
    if (!name) {
        alert('Введите ФИО руководителя');
        return;
    }
    if (supervisors.includes(name)) {
        alert('Этот руководитель уже есть в списке');
        return;
    }
    supervisors.push(name);
    renderSupervisorsList();
    closeAddSupervisorModal();
}

// Удаление руководителя
function deleteSupervisor() {
    if (!selectedSupervisor) {
        alert('Выберите руководителя для удаления');
        return;
    }
    
    if (confirm(`Удалить руководителя "${selectedSupervisor}"?`)) {
        const index = supervisors.indexOf(selectedSupervisor);
        if (index !== -1) {
            supervisors.splice(index, 1);
            selectedSupervisor = null;
            renderSupervisorsList();
        }
    }
}

// Запускаем приложение
document.addEventListener('DOMContentLoaded', init);