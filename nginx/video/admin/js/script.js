// Функция проверки авторизации
async function checkAuth() {
    try {
        const response = await fetch('/admin/api/get/institutions', {
            method: 'GET',
            credentials: 'include' 
        });

        if (!response.ok) {
            if (response.status === 401) {
                setTimeout(() => window.location.href = '/auth/', 2000);
            }
            throw new Error(`HTTP error! status: ${response.status}`);
        }
    } catch (error) {
        showCustomAlert("Введен неправильный логин или пароль");
        console.error('Ошибка проверки авторизации:', error);
        setTimeout(() => window.location.href = '/auth/', 6000);
    }
}

// Класс для работы с API
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
            try {
                const errorData = await response.json();
                console.error('Детали ошибки:', errorData);
                throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
            } catch (e) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
        }
        
        if (response.status === 204) {
            return null;
        }
        
        try {
            return await response.json();
        } catch (e) {
            return null;
        }
    }

    // Institution methods
    async getInstitutions() {
        return this.fetchWithAuth(`${this.adminBaseUrl}/get/institutions`);
    }

    async getInstitutionByINN(inn) {
        return this.fetchWithAuth(`/admin/get/institution?inn=${inn}`);
    }

    async addInstitution(data) {
        const requestData = {
            name: String(data.Name).trim(),
            inn: Number(data.INN),
            columns: Array.isArray(data.Columns) ? data.Columns : [String]
        };

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

    async deleteInstitution(id) {
        if (!id) {
            throw new Error('ID ВУЗа обязательно');
        }

        await this.fetchWithAuth(`${this.adminBaseUrl}/delete/institution`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id })
        });

        return true;
    }

    // Form methods
    async getFormColumns(institutionId) {
        const id = Number(institutionId);
        if (isNaN(id)) throw new Error('ID учреждения должен быть числом');

        return this.fetchWithAuth(`${this.adminBaseUrl}/get/form/columns`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ institution_id: id }) 
        });
    }

    async updateFormColumns(institutionId, columns) {
        const id = Number(institutionId);
        if (isNaN(id)) throw new Error('ID учреждения должен быть числом');

        return this.fetchWithAuth(`${this.adminBaseUrl}/put/institution/columns`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ 
                id: id,          
                columns: columns 
            })
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

// Создание экземпляра API сервиса
const apiService = new ApiService("/admin/api/"); 

// Глобальные переменные
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
    await checkAuth();
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
    }
});

// Загрузка списка вузов
async function loadInstitutions() {
    institutions = await apiService.getInstitutions();
    renderInstitutionsList();
}

// Функции для работы со списком вузов
function renderInstitutionsList() {
    institutionsList.innerHTML = '';
    
    if (!institutions || institutions.length === 0) {
        institutionsList.innerHTML = '<div class="empty">Нет данных об учебных заведениях</div>';
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
    document.querySelectorAll('.list-item').forEach(item => {
        item.classList.remove('selected');
    });
    
    const items = Array.from(document.querySelectorAll('.list-item'));
    const selectedItem = items.find(item => item.textContent === institution.name);
    if (selectedItem) {
        selectedItem.classList.add('selected');
    }
    
    selectedInstitution = institution;
    actionPanel.classList.remove('hidden');

    fetchTemplate();
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
        showCustomAlert('Заполните все обязательные поля');
        return;
    }

    try {
        await apiService.addInstitution({
            Name: name,
            INN: inn,
            Columns: [String] 
        });
        showCustomAlert('Учебное заведение успешно добавлено!');
        closeAddInstitutionModal();
        await loadInstitutions(); 
    } catch (error) {
        console.error('Ошибка:', error);
        showCustomAlert(`Ошибка: ${error.message}`);
    }
}
// Функции для редактирования форм
// Функция открытия модального окна редактора формы
async function openFormEditorModal() {
    if (!selectedInstitution) return;
    
    try {
        currentInstitutionName.textContent = selectedInstitution.name;
        
        await renderFormFields(selectedInstitution);
        
        formEditorModal.style.display = 'block';
    } catch (error) {
        console.error('Ошибка при загрузке формы:', error);
        showCustomAlert('Произошла ошибка при загрузке формы');
    }
}

// Функция закрытия модального окна редактора формы
function closeFormEditorModal() {
    formEditorModal.style.display = 'none';
}

// Функция отображения полей формы
async function renderFormFields(institution) {
    formFieldsContainer.innerHTML = '<div class="loader">Загрузка полей формы...</div>';

    try {
        // Запрашиваем с сервера данные о полях формы для конкретного учреждения
        const response = await apiService.getFormColumns(institution.id);
        console.log("Сервер вернул поля формы:", response);

        const rawFields = Array.isArray(response)
            ? response 
            : (response?.columns || response?.data?.columns || []); 

        if (rawFields.length === 0) {
            formFieldsContainer.innerHTML = `
                <div class="empty-form">
                    <p>Нет настроенных полей</p>
                </div>
            `;
            return;
        }

        // Очищаем контейнер перед добавлением новых полей
        formFieldsContainer.innerHTML = '';

        // Для каждого поля создаем элемент интерфейса
        rawFields.forEach((field, index) => {
            let fieldData = {};

            if (typeof field === 'string') {
                fieldData = {
                    label: field,       
                    type: 'text',       
                    required: false     
                };
            } 

            else {
                fieldData = {
                    label: field.label || '',           
                    type: field.type || 'text',         
                    required: field.required || false  
                };
            }

            // Создаем DOM-элемент для поля формы
            const fieldElement = document.createElement('div');
            fieldElement.className = 'form-field';
            
            // Заполняем HTML-структуру поля
            fieldElement.innerHTML = `
                <div class="form-group">
                    <label>Название поля:</label>
                    <input type="text" value="${fieldData.label}" class="field-label">
                </div>
                <button class="remove-field-btn">Удалить</button>
            `;
            
            // Добавляем поле в контейнер
            formFieldsContainer.appendChild(fieldElement);
        });
    } catch (error) {
        console.error('Ошибка при загрузке полей формы:', error);
        formFieldsContainer.innerHTML = `
            <div class="error">
                Ошибка загрузки: ${error.message}
                <button onclick="renderFormFields(selectedInstitution)">Повторить</button>
            </div>
        `;
    }
}

// Функция добавления нового поля в форму
function addFormField() {
    // Создаем контейнер для нового поля
    const fieldElement = document.createElement('div');
    fieldElement.className = 'form-field';
    
    // Заполняем стандартную структуру нового поля
    fieldElement.innerHTML = `
        <div class="form-group">
            <label>Название поля:</label>
            <input type="text" class="field-label" placeholder="Введите название">
        </div>
        <button class="remove-field-btn">Удалить</button>
    `;
    
    // Добавляем поле в контейнер формы
    formFieldsContainer.appendChild(fieldElement);
}

// Функция сохранения изменений формы
async function saveFormFields() {
    // Проверяем, выбрано ли учебное заведение
    if (!selectedInstitution) return;

    const fields = []; // Массив для хранения данных полей
    let hasEmptyField = false; // Флаг наличия пустых полей

    // Собираем данные всех полей формы
    document.querySelectorAll('.form-field').forEach(fieldEl => {
        // Получаем значение названия поля (убираем пробелы по краям)
        const label = fieldEl.querySelector('.field-label').value.trim();

        // Проверяем, заполнено ли поле
        if (!label) {
            hasEmptyField = true;
        } else {
            // Добавляем название поля в массив
            fields.push(label);
        }
    });

    // Проверяем наличие незаполненных полей
    if (hasEmptyField) {
        showCustomAlert('Пожалуйста, заполните все поля или удалите пустые перед сохранением.');
        return;
    }

    if (fields.length === 0) {
        alert('Невозможно сохранить: форма пуста.');
        return;
    }

    console.log("Отправляемые поля:", fields);

    try {
        await apiService.updateFormColumns(selectedInstitution.id, fields);
        
        closeFormEditorModal();
        
        showCustomAlert('Форма успешно сохранена!');
    } catch (error) {
        console.error('Ошибка при сохранении формы:', error);
        showCustomAlert('Произошла ошибка при сохранении формы');
    }
}

// Обработчик выхода
document.getElementById('logoutBtn')?.addEventListener('click', async () => {
    try {
        await fetch('/auth/api/logout', { 
            method: 'POST',
            credentials: 'include' 
        });
        window.location.href = '/auth/';
    } catch (error) {
        console.error('Ошибка выхода:', error);
    }
});

// Обработчик удаления поля формы
formFieldsContainer.addEventListener('click', (e) => {
    if (e.target.classList.contains('remove-field-btn')) {
        e.target.closest('.form-field').remove();
    }
});

// Функция удаления вуза
async function deleteInstitution() {
    if (!selectedInstitution) {
        showCustomAlert('Выберите учебное заведение для удаления');
        return;
    }
    
    if (confirm(`Удалить учебное заведение "${selectedInstitution.name}"?`)) {
        try {
            await apiService.deleteInstitution(selectedInstitution.id);
            institutions = institutions.filter(m => m.id !== selectedInstitution.id);
            selectedInstitution = null;
            renderInstitutionsList();
            actionPanel.classList.add('hidden'); 
            showCustomAlert('Учебное заведение успешно удалено');
        } catch (error) {
            console.error('Ошибка удаления учебного заведения:', error);
            showCustomAlert(`Ошибка удаления: ${error.message}`);
        }
    }
}

// Назначение обработчиков
editFormBtn.addEventListener('click', openFormEditorModal);
addFieldBtn.addEventListener('click', addFormField);
saveFormBtn.addEventListener('click', saveFormFields);
cancelFormBtn.addEventListener('click', closeFormEditorModal);
closeFormBtn.addEventListener('click', closeFormEditorModal);

// Обработчик загрузки шаблона
document.getElementById('uploadTemplateBtn').addEventListener('click', () => {
    if (!selectedInstitution) {
        showCustomAlert("Сначала выберите учебное заведение.");
        return;
    }
    document.getElementById('templateUploadInput').click();
});

document.getElementById('templateUploadInput').addEventListener('change', function () {
    const file = this.files[0];
    if (!file || !file.name.endsWith('.docx')) {
        showCustomAlert('Пожалуйста, выберите файл формата .docx');
        return;
    }

    if (!selectedInstitution) {
        showCustomAlert("Сначала выберите учебное заведение.");
        return;
    }

    const formData = new FormData();
    formData.append('institution_id', selectedInstitution.id);
    formData.append('file', file);

    fetch('http://generator/api/templates/upload', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => { throw new Error(text || 'Ошибка при загрузке шаблона'); });
        }
        return response.json();
    })
    .then(result => {
        showCustomAlert('Шаблон успешно загружен для "' + selectedInstitution.name + '"');
    })
    .catch(error => {
        showCustomAlert('Ошибка: ' + error.message);
    });
});

// Обработчик генерации документа
generateDocBtn.addEventListener('click', async () => {
    if (!selectedInstitution) {
        showCustomAlert("Выберите учебное заведение перед генерацией документа.");
        return;
    }

    try {
        const response = await fetch("http://generator/api/documents/generate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                institution_id: selectedInstitution.id
            })
        });

        if (!response.ok) {
            const errorData = await response.text();
            throw new Error(errorData || "Ошибка генерации документа.");
        }

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);

        const a = document.createElement("a");
        a.href = url;
        a.download = "generated.docx";
        document.body.appendChild(a);
        a.click();

        a.remove();
        window.URL.revokeObjectURL(url);
    } catch (error) {
        console.error("Ошибка генерации:", error);
        showCustomAlert("Ошибка генерации: " + error.message);
    }
});

// Функции для кастомных уведомлений
function showCustomAlert(message) {
    document.getElementById('customAlertMessage').textContent = message;
    document.getElementById('customAlert').style.display = 'block';
    document.getElementById('customAlertOverlay').style.display = 'block';
}

function hideCustomAlert() {
    document.getElementById('customAlert').style.display = 'none';
    document.getElementById('customAlertOverlay').style.display = 'none';
}