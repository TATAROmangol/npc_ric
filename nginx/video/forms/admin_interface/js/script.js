// Данные
let universities = {
    "МГУ": {
        details: {
            inn: "1234567890",
            address: "Москва, Ленинские горы, 1",
            phone: "+7 (495) 939-1000"
        },
        formFields: [
            { name: "full_name", label: "ФИО студента", type: "text", required: true },
            { name: "practice_start", label: "Начало практики", type: "date", required: true },
            { name: "practice_end", label: "Конец практики", type: "date", required: true }
        ]
    },
    "ВШЭ": {
        details: {
            inn: "0987654321",
            address: "Москва, Покровский бульвар, 11",
            phone: "+7 (495) 771-3232"
        },
        formFields: [
            { name: "student_name", label: "ФИО студента", type: "text", required: true },
            { name: "discipline", label: "Дисциплина", type: "text", required: false }
        ]
    }
};


let selectedUniversity = null;

// DOM элементы
const institutionsList = document.getElementById('institutionsList');
const actionPanel = document.getElementById('actionPanel');
const openModalBtn = document.getElementById('openModalBtn');
const generateDocBtn = document.getElementById('generateDocBtn');
const editFormBtn = document.getElementById('editFormBtn');
const deleteUniversityBtn = document.getElementById('deleteUniversityBtn');

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
const currentUniversityName = document.getElementById('currentUniversityName');
const formFieldsContainer = document.getElementById('formFieldsContainer');
const addFieldBtn = document.getElementById('addFieldBtn');
const saveFormBtn = document.getElementById('saveFormBtn');
const cancelFormBtn = document.getElementById('cancelFormBtn');
const closeFormBtn = document.querySelector('.close-form-btn');

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', () => {
    renderUniversitiesList();
    

    // Обработчики для модального окна добавления вуза
    openModalBtn.addEventListener('click', openAddUniversityModal);
    submitBtn.addEventListener('click', addUniversity);
    cancelBtn.addEventListener('click', closeAddUniversityModal);
    closeBtn.addEventListener('click', closeAddUniversityModal);
    
    // Обработчик удаления вуза
    deleteUniversityBtn.addEventListener('click', deleteUniversity);

    // Закрытие модальных окон при клике вне окна
    window.addEventListener('click', (e) => {
        if (e.target === addInstitutionModal) closeAddUniversityModal();
        if (e.target === formEditorModal) closeFormEditorModal();
    });
});

// Функции для работы со списком вузов
function renderUniversitiesList() {
    institutionsList.innerHTML = '';
    Object.keys(universities).forEach(universityName => {
        const item = document.createElement('div');
        item.className = 'list-item';
        item.textContent = universityName;
        item.addEventListener('click', () => selectUniversity(universityName));
        institutionsList.appendChild(item);
    });
}


function selectUniversity(universityName) {
    // Снимаем выделение со всех элементов
    document.querySelectorAll('.list-item').forEach(item => {
        item.classList.remove('selected');
    });
    
    // Выделяем выбранный
    const items = Array.from(document.querySelectorAll('.list-item'));
    const selectedItem = items.find(item => item.textContent === universityName);
    if (selectedItem) {
        selectedItem.classList.add('selected');
    }
    
    selectedUniversity = universityName;
    actionPanel.classList.remove('hidden');
}



// Функции для модального окна добавления вуза
function openAddUniversityModal() {
    institutionNameInput.value = '';
    institutionINNInput.value = '';
    institutionAddressInput.value = '';
    institutionNumberInput.value = '';
    addInstitutionModal.style.display = 'block';
}

function closeAddUniversityModal() {
    addInstitutionModal.style.display = 'none';
}


function addUniversity() {
    const name = institutionNameInput.value.trim();
    const inn = institutionINNInput.value.trim();
    const address = institutionAddressInput.value.trim();
    const phone = institutionNumberInput.value.trim();
    
    if (!name) {
        alert('Введите название университета');
        return;
    }
    
    if (universities[name]) {
        alert('Такой университет уже существует');
        return;
    }
    
    // Добавление нового университета
    universities[name] = {
        details: { inn, address, phone },
        formFields: [] // Начальный пустой список полей формы
    };
    
    renderUniversitiesList();
    closeAddUniversityModal();
    alert('Университет успешно добавлен!');
}


// Функции для модального окна редактирования формы
function openFormEditorModal() {
    if (!selectedUniversity) return;
    
    currentUniversityName.textContent = selectedUniversity;
    renderFormFields(selectedUniversity);
    formEditorModal.style.display = 'block';
}

function closeFormEditorModal() {
    formEditorModal.style.display = 'none';
}

function renderFormFields(universityName) {
    formFieldsContainer.innerHTML = '';
    
    const fields = universities[universityName].formFields || [];
    
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
            <button>Удалить</button>
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

function saveFormFields() {
    if (!selectedUniversity) return;
    
    const fields = [];
    document.querySelectorAll('.form-field').forEach(fieldEl => {
        fields.push({
            label: fieldEl.querySelector('.field-label').value,
            type: fieldEl.querySelector('.field-type').value,
            required: fieldEl.querySelector('.field-required').checked,
            name: fieldEl.querySelector('.field-label').value.toLowerCase().replace(/\s+/g, '_')
        });
    });
    
    universities[selectedUniversity].formFields = fields;
    closeFormEditorModal();
    alert('Форма сохранена!');
}

// Обработчик удаления поля формы
formFieldsContainer.addEventListener('click', (e) => {
    if (e.target.classList.contains('remove-field-btn')) {
        e.target.closest('.form-field').remove();
    }
});

// Функция удаления вуза
function deleteUniversity() {
    if (!selectedUniversity) return;
    
    if (confirm(`Удалить университет "${selectedUniversity}"?`)) {
        delete universities[selectedUniversity];
        selectedUniversity = null;
        actionPanel.classList.add('hidden');
        renderUniversitiesList();
    }
}

editFormBtn.addEventListener('click', openFormEditorModal);
addFieldBtn.addEventListener('click', addFormField);
saveFormBtn.addEventListener('click', saveFormFields);
