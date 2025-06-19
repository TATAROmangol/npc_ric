// Объект для хранения данных о вузах и их формах
let universitiesData = {};

// DOM элементы
const universitySelect = document.getElementById('university-select');
const supervisorSelect = document.getElementById('supervisor-select');
const formFieldsContainer = document.getElementById('form-fields-container');
const submitBtn = document.getElementById('submit-btn');
const dynamicForm = document.getElementById('dynamic-form');

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', async () => {
    try {
        await Promise.all([
            loadUniversities(),
            loadSupervisors()
        ]);
        
        // Инициализация Flatpickr для полей даты
        initDatePickers();
        
        // Добавляем обработчики событий
        setupEventListeners();
        
    } catch (error) {
        console.error('Ошибка инициализации:', error);
        showCustomAlert('Произошла ошибка при загрузке данных');
    }
});

// Загрузка списка университетов
async function loadUniversities() {
    try {
        const response = await fetch('/forms/api/get/institutions');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        
        // Формируем структуру данных
        universitiesData = data.reduce((acc, university) => {
            acc[university.name] = {
                id: university.id,
                fields: [] 
            };
            return acc;
        }, {});

        // Заполняем выпадающий список
        populateSelect(universitySelect, Object.keys(universitiesData), 'Выберите учебное заведение');
        
    } catch (error) {
        console.error('Ошибка загрузки университетов:', error);
        throw error;
    }
}

// Загрузка списка руководителей
async function loadSupervisors() {
    try {
        const response = await fetch('/forms/api/get/mentors');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const mentors = await response.json();
        const mentorNames = mentors.map(mentor => mentor.name);
        
        populateSelect(supervisorSelect, mentorNames, 'Выберите руководителя');
        
    } catch (error) {
        console.error('Ошибка загрузки руководителей:', error);
        throw error;
    }
}

// Универсальная функция заполнения select элемента
function populateSelect(selectElement, options, placeholderText) {
    selectElement.innerHTML = '';
    
    // Добавляем placeholder
    const placeholderOption = document.createElement('option');
    placeholderOption.value = '';
    placeholderOption.textContent = placeholderText;
    placeholderOption.selected = true;
    placeholderOption.disabled = true;
    selectElement.appendChild(placeholderOption);
    
    // Добавляем варианты выбора
    options.forEach(option => {
        const optionElement = document.createElement('option');
        optionElement.value = option;
        optionElement.textContent = option;
        selectElement.appendChild(optionElement);
    });
}

// Настройка обработчиков событий
function setupEventListeners() {
    // Обработчик изменения выбранного университета
    universitySelect.addEventListener('change', async function() {
        const selectedUniName = this.value;

        resetForm();
        
        if (selectedUniName && universitiesData[selectedUniName]) {
            try {
                // Загружаем поля формы, если еще не загружены
                if (universitiesData[selectedUniName].fields.length === 0) {
                    await loadFormFields(selectedUniName);
                }
                
                // Создаем поля формы
                createFormFields(universitiesData[selectedUniName].fields);
                submitBtn.classList.remove('hidden');
                
            } catch (error) {
                console.error('Ошибка загрузки полей формы:', error);
                showCustomAlert('Не удалось загрузить поля формы');
            }
        }
    });
    
    // Обработчик отправки формы
    dynamicForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        try {
            const isValid = await validateAndSubmitForm();
            if (isValid) {
                resetForm();
            }
        } catch (error) {
            console.error('Ошибка отправки формы:', error);
            showCustomAlert(`Ошибка: ${error.message}`);
        }
    });

    formFieldsContainer.addEventListener('input', clearFieldError);
}

// Загрузка полей формы для выбранного университета
async function loadFormFields(universityName) {
    try {
        const universityId = universitiesData[universityName].id;
        
        const response = await fetch('/forms/api/get/form/columns', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ institution_id: universityId })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const fields = await response.json();
        
        // Нормализуем данные полей
        universitiesData[universityName].fields = Array.isArray(fields)
            ? fields.map(field => ({
                name: field.toLowerCase().replace(/\s+/g, '_'),
                label: field,
                type: 'text',
                required: true
            }))
            : fields;
            
    } catch (error) {
        console.error('Ошибка загрузки полей формы:', error);
        throw error;
    }
}

// Создание полей формы
function createFormFields(fieldsConfig) {
    formFieldsContainer.innerHTML = '';
    
    fieldsConfig.forEach(field => {
        const fieldGroup = document.createElement('div');
        fieldGroup.className = 'form-group';
        
        // Создаем label
        const label = document.createElement('label');
        label.htmlFor = field.name;
        label.textContent = field.label;
        
        // Создаем input/select в зависимости от типа
        let input;
        if (field.type === 'select') {
            input = document.createElement('select');
            input.id = field.name;
            input.name = field.name;
            input.className = 'form-control';
            
            // Добавляем опции, если они есть
            if (field.options) {
                field.options.forEach(option => {
                    const optionElement = document.createElement('option');
                    optionElement.value = option.value;
                    optionElement.textContent = option.label;
                    input.appendChild(optionElement);
                });
            }
        } else {
            input = document.createElement('input');
            input.id = field.name;
            input.name = field.name;
            input.type = field.type || 'text';
            input.className = 'form-control';
            
            // Особые классы для специальных типов
            if (input.type === 'date') {
                input.classList.add('datepicker');
            }
        }
        
        fieldGroup.appendChild(label);
        fieldGroup.appendChild(input);
        formFieldsContainer.appendChild(fieldGroup);
    });
    
    // Инициализируем datepicker для полей даты
    initDatePickers();
}

// Инициализация Flatpickr для полей даты
function initDatePickers() {
    document.querySelectorAll('.datepicker').forEach(el => {
        flatpickr(el, {
            dateFormat: "d.m.Y",
            locale: "ru",
            allowInput: true
        });
    });
}

// Валидация и отправка формы
async function validateAndSubmitForm() {
    const selectedUniName = universitySelect.value;
    
    if (!selectedUniName) {
        showCustomAlert('Выберите учебное заведение');
        return false;
    }
    
    if (!supervisorSelect.value) {
        markFieldAsInvalid(supervisorSelect, 'Выберите руководителя');
        return false;
    }
    
    const university = universitiesData[selectedUniName];
    let isValid = true;
    
    // Проверяем все обязательные поля
    university.fields.forEach(field => {
        const input = document.getElementById(field.name);
        const value = input.value.trim();
        
        if (field.required && !value) {
            markFieldAsInvalid(input, 'Это поле обязательно для заполнения');
            isValid = false;
        } else if (value && !validateFieldValue(value, field.type)) {
            markFieldAsInvalid(input, getValidationMessage(field.type));
            isValid = false;
        }
    });
    
    if (!isValid) {
        showCustomAlert('Пожалуйста, заполните все обязательные поля корректно');
        return false;
    }
    
    // Формируем данные для отправки
    const formData = {
        institution_id: university.id,
        supervisor: supervisorSelect.value,
        fields: university.fields.map(field => ({
            name: field.name,
            value: document.getElementById(field.name).value.trim()
        }))
    };
    
    // Отправляем данные
    try {
        const token = localStorage.getItem('authToken');
        const headers = {
            'Content-Type': 'application/json'
        };
        
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        const response = await fetch('/forms/api/post/form', {
            method: 'POST',
            headers,
            body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Ошибка сервера');
        }
        
        showCustomAlert('Форма успешно отправлена!');
        return true;
        
    } catch (error) {
        console.error('Ошибка отправки формы:', error);
        throw error;
    }
}

// Валидация значения поля по типу
function validateFieldValue(value, type) {
    if (!value) return true;
    
    switch (type) {
        case 'number':
            return !isNaN(value) && isFinite(value);
        case 'email':
            return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
        case 'date':
            return /^\d{2}\.\d{2}\.\d{4}$/.test(value);
        default:
            return true;
    }
}

// Получение сообщения об ошибке для типа поля
function getValidationMessage(type) {
    switch (type) {
        case 'number': return 'Введите корректное число';
        case 'email': return 'Введите корректный email';
        case 'date': return 'Введите дату в формате ДД.ММ.ГГГГ';
        default: return 'Некорректное значение';
    }
}

// Помечаем поле как невалидное
function markFieldAsInvalid(element, message) {
    element.classList.add('invalid');
    
    let errorElement = element.nextElementSibling;
    if (!errorElement || !errorElement.classList.contains('error-message')) {
        errorElement = document.createElement('div');
        errorElement.className = 'error-message';
        element.parentNode.insertBefore(errorElement, element.nextSibling);
    }
    
    errorElement.textContent = message;
}

// Очищаем ошибку поля
function clearFieldError(e) {
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'SELECT') {
        e.target.classList.remove('invalid');
        
        const errorElement = e.target.nextElementSibling;
        if (errorElement && errorElement.classList.contains('error-message')) {
            errorElement.remove();
        }
    }
}

// Сброс формы
function resetForm() {
    formFieldsContainer.innerHTML = '';
    submitBtn.classList.add('hidden');
    universitySelect.classList.remove('invalid');
    
    const errorElement = supervisorSelect.nextElementSibling;
    if (errorElement && errorElement.classList.contains('error-message')) {
        errorElement.remove();
    }
}

// Функции для показа/скрытия кастомных уведомлений
function showCustomAlert(message) {
    const alertElement = document.getElementById('customAlert');
    const messageElement = document.getElementById('customAlertMessage');
    
    messageElement.textContent = message;
    alertElement.style.display = 'block';
    
    // Автоматическое скрытие через 5 секунд
    setTimeout(() => {
        alertElement.style.display = 'none';
    }, 5000);
}
