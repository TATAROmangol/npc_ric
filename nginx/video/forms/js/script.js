let universitiesData = {};

// DOM элементы
const universitySelect = document.getElementById('university-select');
const supervisorSelect = document.getElementById('supervisor-select');
const formFieldsContainer = document.getElementById('form-fields-container');
const submitBtn = document.getElementById('submit-btn');
const dynamicForm = document.getElementById('dynamic-form');

// Загрузка университетов
async function loadUniversities() {
    try {
        const response = await fetch('/forms/api/get/institutions');
        if (!response.ok) throw new Error('Ошибка загрузки университетов');

        const data = await response.json();

        universitiesData = {};
        data.forEach(university => {
            universitiesData[university.name] = {
                id: university.id,
                fields: [] // загружаем позже
            };
        });

        populateUniversitySelect();
    } catch (error) {
        console.error('Ошибка при загрузке учебных заведений:', error);
        alert('Не удалось загрузить список учебных заведений');
    }
}

// Заполнение списка университетов
function populateUniversitySelect() {
    universitySelect.innerHTML = '';

    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'Выберите учебное заведение';
    defaultOption.selected = true;
    defaultOption.disabled = true;
    universitySelect.appendChild(defaultOption);

    Object.keys(universitiesData).forEach(name => {
        const option = document.createElement('option');
        option.value = name;
        option.textContent = name;
        universitySelect.appendChild(option);
    });
}

// Заполнение списка руководителей
async function populateSupervisorSelect() {
    supervisorSelect.innerHTML = '';

    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'Выберите руководителя';
    defaultOption.selected = true;
    defaultOption.disabled = true;
    supervisorSelect.appendChild(defaultOption);

    try {
        const response = await fetch('/forms/api/get/mentors');
        if (!response.ok) {
            throw new Error('Ошибка загрузки руководителей');
        }

        const mentors = await response.json();

        mentors.forEach(mentor => {
            const option = document.createElement('option');
            option.value = mentor.name;
            option.textContent = mentor.name;
            supervisorSelect.appendChild(option);
        });
    } catch (error) {
        console.error('Ошибка при загрузке руководителей:', error);
        alert('Не удалось загрузить список руководителей');
    }
}

// Обработка смены университета
universitySelect.addEventListener('change', async function () {
    const selectedUniversity = this.value;
    formFieldsContainer.innerHTML = '';
    submitBtn.classList.add('hidden');

    if (selectedUniversity && universitiesData[selectedUniversity]) {
        if (universitiesData[selectedUniversity].fields.length === 0) {
            try {
                const response = await fetch('/forms/api/get/form/columns', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        institution_id: universitiesData[selectedUniversity].id
                    })
                });

                if (!response.ok) throw new Error('Ошибка загрузки полей формы');
                const fields = await response.json();

                const parsedFields = Array.isArray(fields)
                    ? fields.map(name => ({
                        name: name.toLowerCase().replace(/\s+/g, '_'),
                        label: name,
                        type: 'text',
                        required: false
                    }))
                    : fields;

                universitiesData[selectedUniversity].fields = parsedFields;
            } catch (error) {
                console.error('Ошибка при загрузке полей формы:', error);
                alert('Не удалось загрузить поля формы');
                return;
            }
        }

        createFormFields(universitiesData[selectedUniversity].fields);
        submitBtn.classList.remove('hidden');
    }
});

// Создание полей формы
function createFormFields(fieldsConfig) {
    formFieldsContainer.innerHTML = '';

    fieldsConfig.forEach(field => {
        const fieldGroup = document.createElement('div');
        fieldGroup.className = 'form-group';

        const label = document.createElement('label');
        label.htmlFor = field.name;
        label.textContent = field.label + (field.required ? '*' : '');

        let input;
        
        if (field.type === 'select') {
            input = document.createElement('select');
            input.id = field.name;
            input.name = field.name;
            input.required = !!field.required;
            input.className = 'form-control';
            
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
            input.required = !!field.required;
            input.className = 'form-control';
        }

        if (input.type === 'date') {
            input.classList.add('datepicker');
        }

        fieldGroup.appendChild(label);
        fieldGroup.appendChild(input);
        formFieldsContainer.appendChild(fieldGroup);
    });

    document.querySelectorAll('.datepicker').forEach(el => {
        flatpickr(el, {
            dateFormat: "d.m.Y",
            locale: "ru"
        });
    });
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

function checkRequiredFields(fieldsConfig) {
    let allRequiredFilled = true;
    
    fieldsConfig.forEach(field => {
        if (field.required) {
            const input = document.getElementById(field.name);
            const value = input.value.trim();
            
            if (!value) {
                input.classList.add('error');
                
                let errorElement = input.nextElementSibling;
                if (!errorElement || !errorElement.classList.contains('error-message')) {
                    errorElement = document.createElement('div');
                    errorElement.className = 'error-message';
                    errorElement.textContent = 'Это поле обязательно для заполнения';
                    input.parentNode.insertBefore(errorElement, input.nextSibling);
                }
                
                allRequiredFilled = false;
            } else {
                input.classList.remove('error');
                const errorElement = input.nextElementSibling;
                if (errorElement && errorElement.classList.contains('error-message')) {
                    errorElement.remove();
                }
            }
        }
    });
    
    return allRequiredFilled;
}

// Валидация всей формы перед отправкой
function validateForm(fieldsConfig) {
    let isValid = true;
    
    fieldsConfig.forEach(field => {
        const input = document.getElementById(field.name);
        const value = input.value.trim();
        
        // Для обязательных полей проверка уже была, проверяем только заполненные
        if (value) {
            if (!validateFieldValue(value, field.type)) {
                input.classList.add('error');
                isValid = false;
                
                let errorElement = input.nextElementSibling;
                if (!errorElement || !errorElement.classList.contains('error-message')) {
                    errorElement = document.createElement('div');
                    errorElement.className = 'error-message';
                    input.parentNode.insertBefore(errorElement, input.nextSibling);
                }
                
                errorElement.textContent = getValidationErrorMessage(field.type);
            }
        }
    });
    
    return isValid;
}

function getValidationErrorMessage(type) {
    switch (type) {
        case 'number': return 'Введите корректное число';
        case 'email': return 'Введите корректный email';
        case 'date': return 'Введите дату в формате ДД.ММ.ГГГГ';
        default: return 'Некорректное значение';
    }
}

// Отправка формы
dynamicForm.addEventListener('submit', async function (e) {
    e.preventDefault(); 

    const selectedUniversity = universitySelect.value;
    if (!selectedUniversity) {
        alert('Пожалуйста, выберите учебное заведение');
        return;
    }

    const university = universitiesData[selectedUniversity];
    
    // 1. Проверяем заполнение обязательных полей
    const areRequiredFieldsFilled = checkRequiredFields(university.fields);
    if (!areRequiredFieldsFilled) {
        alert('Пожалуйста, заполните все обязательные поля');
        return;
    }
    
    // 2. Проверяем валидность данных в полях
    const isFormValid = validateForm(university.fields);
    if (!isFormValid) {
        alert('Пожалуйста, исправьте ошибки в форме перед отправкой');
        return;
    }

    // Если все проверки пройдены, собираем данные
    const infoArray = university.fields.map(field => {
        return document.getElementById(field.name).value;
    });

    const formData = {
        institution_id: university.id,
        info: infoArray
    };

    // Отправка данных
    try {
        const token = localStorage.getItem('authToken');
        const headers = {
            'Content-Type': 'application/json'
        };
        if (token) headers['Authorization'] = `Bearer ${token}`;

        const response = await fetch('/forms/api/post/form', {
            method: 'POST',
            headers,
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            throw new Error(await response.text() || 'Ошибка отправки формы');
        }

        alert('Заявление успешно отправлено!');
        dynamicForm.reset();
        formFieldsContainer.innerHTML = '';
        submitBtn.classList.add('hidden');
        universitySelect.value = '';
        
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка при отправке формы: ' + error.message);
    }
});

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', () => {
    loadUniversities();
    populateSupervisorSelect();
});

// Добавляем обработчики для сброса ошибок при вводе
formFieldsContainer.addEventListener('input', (e) => {
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'SELECT') {
        e.target.classList.remove('error');
        const errorElement = e.target.nextElementSibling;
        if (errorElement && errorElement.classList.contains('error-message')) {
            errorElement.remove();
        }
    }
});
