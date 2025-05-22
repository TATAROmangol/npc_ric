class ApiService {
    constructor() {
        this.adminBaseUrl = '/admin/api';
        this.formsBaseUrl = '/forms/api';
    }

    async fetchWithAuth(url, options = {}) {
        const token = localStorage.getItem('authToken');
        if (token) {
            options.headers = {
                ...options.headers,
                'Authorization': `Bearer ${token}`
            };
        }

        const response = await fetch(url, options);
        if (!response.ok) {
            const errorData = await response.json(); // Пытаемся прочитать JSON-ошибку
            console.error('Детали ошибки:', errorData);
            throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
        }
        return response.json();
    }

    // Institution methods
    async getInstitutions() {
        return this.fetchWithAuth(`${this.adminBaseUrl}/get/institutions`);
    }

    async getInstitutionByINN(inn) {
        return this.fetchWithAuth(`/admin/get/institution?inn=${inn}`);
    }

    async addInstitution(data) {
        // Гарантируем правильный формат данных
        const requestData = {
            name: String(data.Name).trim(),
            inn: Number(data.INN),
            columns: Array.isArray(data.Columns) ? data.Columns : [String]
        };

        // Валидация
        if (!requestData.name) throw new Error('Название обязательно');
        if (isNaN(requestData.inn)) throw new Error('ИНН должен быть числом');

        console.log('Отправляемые данные:', requestData); // Для отладки

        return this.fetchWithAuth(`${this.adminBaseUrl}/post/institution`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify(requestData)
        });
    }

    async updateInstitution(data) {
        return this.fetchWithAuth(`${this.adminBaseUrl}/institution`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
    }

    async deleteInstitution(inn) {
        // Преобразуем ИНН в число
        const innNumber = Number(inn);
        if (isNaN(innNumber)) {
            throw new Error('ИНН должен быть числом');
        }

        return this.fetchWithAuth(`${this.adminBaseUrl}/delete/institution`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                inn: innNumber  // Отправляем ИНН как число
            })
        });
    }

    // Form methods
    async getFormColumns(institutionId) {
        return this.fetchWithAuth(`${this.formsBaseUrl}/get/form/columns?institutionId=${institutionId}`);
    }

    async updateFormColumns(institutionId, columns) {
        return this.fetchWithAuth('put/form/columns', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ institutionId, columns })
        });
    }

    // Mentor methods
    async getMentors() {
        return this.fetchWithAuth(`${this.adminBaseUrl}/get/mentors`);
    }

    async addMentor(data) {
        return this.fetchWithAuth(`${this.adminBaseUrl}/post/mentor`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
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
        return this.fetchWithAuth(`${this.adminBaseUrl}/delete/mentor`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id })
        });
    }
}

const apiService = new ApiService("/admin/api/"); 

let institutions = []
let selectedInstitution = null;





// DOM элементы
const institutionsList = document.getElementById('institutionsList');
const actionPanel = document.getElementById('actionPanel');
const openModalBtn = document.getElementById('openModalBtn');
const generateDocBtn = document.getElementById('generateDocBtn');
const editFormBtn = document.getElementById('editFormBtn');
const deleteInstitutionBtn = document.getElementById('deleteInstitutionBtn');

// Элементы модального окна добавления вуза
const addInstitutionModal = document.getElementById('addInstitutionModal');
const institutionNameInput = document.getElementById('institutionName');
const institutionINNInput = document.getElementById('institutionINN');
const institutionAddressInput = document.getElementById('institutionAddress');
const institutionNumberInput = document.getElementById('institutionNumber');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');
const closeBtn = document.querySelector('.close-btn');

// Элементы модального окна редактирования формы
const formEditorModal = document.getElementById('formEditorModal');
const currentInstitutionName = document.getElementById('currentInstitutionName');
const formFieldsContainer = document.getElementById('formFieldsContainer');
const addFieldBtn = document.getElementById('addFieldBtn');
const saveFormBtn = document.getElementById('saveFormBtn');
const cancelFormBtn = document.getElementById('cancelFormBtn');
const closeFormBtn = document.querySelector('.close-form-btn');

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', async () => {
    try {
        await loadInstitutions();
        
        // Обработчики для модального окна добавления вуза
        openModalBtn.addEventListener('click', openAddInstitutionModal);
        submitBtn.addEventListener('click', addInstitution);
        cancelBtn.addEventListener('click', closeAddInstitutionModal);
        closeBtn.addEventListener('click', closeAddInstitutionModal);
        
        // Обработчик удаления вуза
        deleteInstitutionBtn.addEventListener('click', deleteInstitution);

        // Закрытие модальных окон при клике вне окна
        window.addEventListener('click', (e) => {
            if (e.target === addInstitutionModal) closeAddInstitutionModal();
            if (e.target === formEditorModal) closeFormEditorModal();
        });

    } catch (error) {
        console.error('Ошибка при загрузке данных:', error);
        alert('Произошла ошибка при загрузке данных');
    }
});

// Загрузка списка вузов с сервера
async function loadInstitutions() {
    institutions = await apiService.getInstitutions();
    renderInstitutionsList();
}

// Функции для работы со списком вузов
function renderInstitutionsList() {
    institutionsList.innerHTML = '';
    
    if (!institutions || institutions.length === 0) {
        institutionsList.innerHTML = '<div class="empty">Нет данных об университетах</div>';
        return;
    }
    
    institutions.forEach(institution => {
        const item = document.createElement('div');
        item.className = 'list-item';
        item.textContent = institution.name;
        item.dataset.id = institution.id;
        item.addEventListener('click', () => selectInstitution(institution));
        institutionsList.appendChild(item);
    });
}

function selectInstitution(institution) {
    // Снимаем выделение со всех элементов
    document.querySelectorAll('.list-item').forEach(item => {
        item.classList.remove('selected');
    });
    
    // Выделяем выбранный
    const items = Array.from(document.querySelectorAll('.list-item'));
    const selectedItem = items.find(item => item.textContent === institution.name);
    if (selectedItem) {
        selectedItem.classList.add('selected');
    }
    
    selectedInstitution = institution;
    actionPanel.classList.remove('hidden');
}

// Функции для модального окна добавления вуза
function openAddInstitutionModal() {
    institutionNameInput.value = '';
    institutionINNInput.value = '';
    addInstitutionModal.style.display = 'block';
}

function closeAddInstitutionModal() {
    addInstitutionModal.style.display = 'none';
}

async function addInstitution() {
    const name = institutionNameInput.value.trim();
    const inn = institutionINNInput.value.trim();

    if (!name || !inn) {
        alert('Заполните все обязательные поля');
        return;
    }

    try {
        await apiService.addInstitution({
            Name: name,
            INN: inn,
            Columns: [String] // Явно указываем пустой массив columns
        });
        alert('Университет успешно добавлен!');
        closeAddInstitutionModal();
        await loadInstitutions(); // Обновляем список
    } catch (error) {
        console.error('Ошибка:', error);
        alert(`Ошибка: ${error.message}`);
    }
}

// Функции для модального окна редактирования формы
async function openFormEditorModal() {
    if (!selectedInstitution) return;
    
    try {
        currentInstitutionName.textContent = selectedInstitution.name;
        await renderFormFields(selectedInstitution);
        formEditorModal.style.display = 'block';
    } catch (error) {
        console.error('Ошибка при загрузке формы:', error);
        alert('Произошла ошибка при загрузке формы');
    }
}

function closeFormEditorModal() {
    formEditorModal.style.display = 'none';
}

async function renderFormFields(institution) {
    formFieldsContainer.innerHTML = '';
    
    try {
        const formColumns = await apiService.getFormColumns(institution.id);
        const fields = formColumns || [];
        
        fields.forEach((field, index) => {
            const fieldElement = document.createElement('div');
            fieldElement.className = 'form-field';
            fieldElement.innerHTML = `
                <div class="form-group">
                    <label>Название поля:</label>
                    <input type="text" value="${field.label || ''}" class="field-label">
                </div>
                <div class="form-group">
                    <label>Тип поля:</label>
                    <select class="field-type">
                        <option value="text" ${field.type === 'text' ? 'selected' : ''}>Текст</option>
                        <option value="date" ${field.type === 'date' ? 'selected' : ''}>Дата</option>
                        <option value="number" ${field.type === 'number' ? 'selected' : ''}>Число</option>
                        <option value="email" ${field.type === 'email' ? 'selected' : ''}>Email</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>
                        <input type="checkbox" class="field-required" ${field.required ? 'checked' : ''}>
                        Обязательное поле
                    </label>
                </div>
                <button class="remove-field-btn">Удалить</button>
            `;
            formFieldsContainer.appendChild(fieldElement);
        });
    } catch (error) {
        console.error('Ошибка при загрузке полей формы:', error);
        alert('Не удалось загрузить поля формы');
    }
}

function addFormField() {
    const fieldElement = document.createElement('div');
    fieldElement.className = 'form-field';
    fieldElement.innerHTML = `
        <div class="form-group">
            <label>Название поля:</label>
            <input type="text" class="field-label" placeholder="Введите название">
        </div>
        <div class="form-group">
            <label>Тип поля:</label>
            <select class="field-type">
                <option value="text">Текст</option>
                <option value="date">Дата</option>
                <option value="number">Число</option>
                <option value="email">Email</option>
            </select>
        </div>
        <div class="form-group">
            <label>
                <input type="checkbox" class="field-required">
                Обязательное поле
            </label>
        </div>
        <button class="remove-field-btn">Удалить</button>
    `;
    formFieldsContainer.appendChild(fieldElement);
}

async function saveFormFields() {
    if (!selectedInstitution) return;
    
    const fields = [];
    document.querySelectorAll('.form-field').forEach(fieldEl => {
        fields.push({
            label: fieldEl.querySelector('.field-label').value,
            type: fieldEl.querySelector('.field-type').value,
            required: fieldEl.querySelector('.field-required').checked,
            name: fieldEl.querySelector('.field-label').value.toLowerCase().replace(/\s+/g, '_')
        });
    });
    
    try {
        await apiService.updateFormColumns(selectedInstitution.id, fields);
        closeFormEditorModal();
        alert('Форма сохранена!');
    } catch (error) {
        console.error('Ошибка при сохранении формы:', error);
        alert('Произошла ошибка при сохранении формы');
    }
}

// Обработчик удаления поля формы
formFieldsContainer.addEventListener('click', (e) => {
    if (e.target.classList.contains('remove-field-btn')) {
        e.target.closest('.form-field').remove();
    }
});

// Функция удаления вуза
async function deleteInstitution() {
    if (!selectedInstitution) {
        alert('Выберите университет для удаления');
        return;
    }

    if (!confirm(`Удалить университет "${selectedInstitution.name}"?`)) {
        return;
    }

    try {
        await apiService.deleteInstitution(selectedInstitution.inn);
        
        // Обновляем список
        institutions = institutions.filter(inst => inst.inn !== selectedInstitution.inn);
        selectedInstitution = null;
        renderInstitutionsList();
        alert('Университет успешно удалён');
    } catch (error) {
        console.error('Ошибка удаления:', error);
        alert(`Ошибка удаления: ${error.message}`);
    }
}

// Назначение обработчиков
editFormBtn.addEventListener('click', openFormEditorModal);
addFieldBtn.addEventListener('click', addFormField);
saveFormBtn.addEventListener('click', saveFormFields);
cancelFormBtn.addEventListener('click', closeFormEditorModal);
closeFormBtn.addEventListener('click', closeFormEditorModal);
