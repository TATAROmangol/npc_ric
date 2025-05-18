class ApiService {
    constructor(baseUrl = 'http://forms:8080/forms/') {
        this.baseUrl = baseUrl;
    }

    async fetchWithAuth(url, options = {}) {
        const token = localStorage.getItem('authToken');
        if (token) {
            options.headers = {
                ...options.headers,
                'Authorization': `Bearer ${token}`
            };
        }
        
        const response = await fetch(`${this.baseUrl}${url}`, options);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    }

    // Institution methods
    async getInstitutions() {
        return this.fetchWithAuth('/admin/get/institutions');
    }

    async getInstitutionByINN(inn) {
        return this.fetchWithAuth(`/admin/get/institution?inn=${inn}`);
    }

    async addInstitution(data) {
        return this.fetchWithAuth('/admin/post/institution', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
    }

    async updateInstitution(data) {
        return this.fetchWithAuth('/admin/put/institution', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
    }

    async deleteInstitution(inn) {
        return this.fetchWithAuth('/admin/delete/institution', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ inn })
        });
    }

    // Form methods
    async getFormColumns() {
        return this.fetchWithAuth('/forms/get/form/columns');
    }

    async updateFormColumns(institutionId, columns) {
        return this.fetchWithAuth('/admin/put/institution/columns', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ institutionId, columns })
        });
    }

    // Mentor methods
    async getMentors() {
        return this.fetchWithAuth('/admin/get/mentors');
    }

    async addMentor(data) {
        return this.fetchWithAuth('/admin/post/mentor', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
    }

    async updateMentor(data) {
        return this.fetchWithAuth('/admin/put/mentor', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
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
}

const apiService = new ApiService("/admin/api/"); 

let institutions = [];
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
    institutions.forEach(institution => {
        const item = document.createElement('div');
        item.className = 'list-item';
        item.textContent = institution.name;
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
    institutionAddressInput.value = '';
    institutionNumberInput.value = '';
    addInstitutionModal.style.display = 'block';
}

function closeAddInstitutionModal() {
    addInstitutionModal.style.display = 'none';
}

async function addInstitution() {
    const name = institutionNameInput.value.trim();
    const inn = institutionINNInput.value.trim();
    const address = institutionAddressInput.value.trim();
    const phone = institutionNumberInput.value.trim();
    
    if (!name || !inn) {
        alert('Название и ИНН университета обязательны для заполнения');
        return;
    }
    
    try {
        const newInstitution = await apiService.addInstitution({
            name,
            inn,
            address,
            phone,
            formFields: [] // Начальный пустой список полей формы
        });
        
        institutions.push(newInstitution);
        renderInstitutionsList();
        closeAddInstitutionModal();
        alert('Университет успешно добавлен!');
    } catch (error) {
        console.error('Ошибка при добавлении университета:', error);
        alert('Произошла ошибка при добавлении университета');
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
    
    // Загружаем поля формы для выбранного учреждения
    const formColumns = await apiService.getFormColumns(institution.id);
    const fields = formColumns || [];
    
    fields.forEach((field, index) => {
        const fieldElement = document.createElement('div');
        fieldElement.className = 'form-field';
        fieldElement.innerHTML = `
            <div class="form-group">
                <label>Название поля:</label>
                <input type="text" value="${field.label}" class="field-label">
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
                    Обязательное поле
                    <input type="checkbox" class="field-required" ${field.required ? 'checked' : ''}>
                </label>
            </div>
            <button class="remove-field-btn">Удалить</button>
        `;
        formFieldsContainer.appendChild(fieldElement);
    });
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
    if (!selectedInstitution) return;
    
    if (confirm(`Удалить университет "${selectedInstitution.name}"?`)) {
        try {
            await apiService.deleteInstitution(selectedInstitution.inn);
            institutions = institutions.filter(inst => inst.inn !== selectedInstitution.inn);
            selectedInstitution = null;
            actionPanel.classList.add('hidden');
            renderInstitutionsList();
            alert('Университет успешно удален');
        } catch (error) {
            console.error('Ошибка при удалении университета:', error);
            alert('Произошла ошибка при удалении университета');
        }
    }
}

// Назначение обработчиков
editFormBtn.addEventListener('click', openFormEditorModal);
addFieldBtn.addEventListener('click', addFormField);
saveFormBtn.addEventListener('click', saveFormFields);
cancelFormBtn.addEventListener('click', closeFormEditorModal);
closeFormBtn.addEventListener('click', closeFormEditorModal);
