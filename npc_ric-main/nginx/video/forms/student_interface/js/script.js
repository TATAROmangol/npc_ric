let universitiesData = {}; 

let supervisorData = {
    "supervisor1": "Иванов Иван Иванович",
    "supervisor2": "Петров Петр Петрович",
    "supervisor3": "Максимов Максим Максимович"
};

// DOM элементы
const universitySelect = document.getElementById('university-select');
const supervisorSelect = document.getElementById('supervisor-select');
const formFieldsContainer = document.getElementById('form-fields-container');
const submitBtn = document.getElementById('submit-btn');
const dynamicForm = document.getElementById('dynamic-form');

// Функция для загрузки данных университетов с сервера
async function loadUniversities() {
    try {
        const response = await fetch('/forms/api/get/institutions');
        if (!response.ok) {
            throw new Error('Ошибка загрузки данных университетов');
        }
        const data = await response.json();
        
        // Преобразуем данные сервера в нужный нам формат
        universitiesData = {};
        data.forEach(university => {
            universitiesData[university.name] = {
                fields: university.formFields || []
            };
        });
        
        // Заполняем список университетов
        populateUniversitySelect();
    } catch (error) {
        console.error('Ошибка при загрузке университетов:', error);
        alert('Не удалось загрузить список университетов');
    }
}

// Функция для заполнения выпадающего списка университетов
function populateUniversitySelect() {
    universitySelect.innerHTML = ''; // Очищаем список
    
    // Добавляем пустой вариант по умолчанию
    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'Выберите университет';
    defaultOption.selected = true;
    defaultOption.disabled = true;
    universitySelect.appendChild(defaultOption);
    
    // Добавляем университеты из загруженных данных
    Object.keys(universitiesData).forEach(universityName => {
        const option = document.createElement('option');
        option.value = universityName;
        option.textContent = universityName;
        universitySelect.appendChild(option);
    });
}

// Заполняем список руководителей практик 
function populateSupervisorSelect() {
    supervisorSelect.innerHTML = ''; 
    
    // Добавляем пустой вариант по умолчанию
    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'Выберите руководителя';
    defaultOption.selected = true;
    defaultOption.disabled = true;
    supervisorSelect.appendChild(defaultOption);
    
    // Добавляем руководителей
    Object.values(supervisorData).forEach(supervisorName => {
        const option = document.createElement('option');
        option.value = supervisorName;
        option.textContent = supervisorName;
        supervisorSelect.appendChild(option);
    });
}

// Обработчик изменения выбора университета
universitySelect.addEventListener('change', async function() {
    const selectedUniversity = this.value;
    
    // Очищаем предыдущие поля
    formFieldsContainer.innerHTML = '';
    submitBtn.classList.add('hidden');
    
    if (selectedUniversity && universitiesData[selectedUniversity]) {
        // Проверяем, есть ли уже поля для этого университета
        if (universitiesData[selectedUniversity].fields.length === 0) {
            // Если нет, загружаем с сервера
            try {
                const response = await fetch(`/forms/api/get/form/columns?institution=${encodeURIComponent(selectedUniversity)}`);
                if (!response.ok) {
                    throw new Error('Ошибка загрузки полей формы');
                }
                const fields = await response.json();
                universitiesData[selectedUniversity].fields = fields;
            } catch (error) {
                console.error('Ошибка при загрузке полей формы:', error);
                alert('Не удалось загрузить поля формы для выбранного университета');
                return;
            }
        }
        
        // Создаем поля формы
        createFormFields(universitiesData[selectedUniversity].fields);
        submitBtn.classList.remove('hidden');
    }
});

// Функция создания полей формы (остается без изменений)
function createFormFields(fieldsConfig) {
    fieldsConfig.forEach(field => {
        const fieldGroup = document.createElement('div');
        fieldGroup.className = 'form-group';
        
        const label = document.createElement('label');
        label.htmlFor = field.name;
        label.textContent = field.label + (field.required ? '*' : '');
        
        let input;
        if (field.type === 'select') {
            input = document.createElement('select');
            // Здесь можно добавить options если нужно
        } else {
            input = document.createElement('input');
            input.type = field.type;
            
            // Добавляем flatpickr для полей даты
            if (field.type === 'date') {
                input.classList.add('datepicker');
            }
        }
        
        input.id = field.name;
        input.name = field.name;
        input.required = field.required;
        input.className = 'form-control';
        
        fieldGroup.appendChild(label);
        fieldGroup.appendChild(input);
        formFieldsContainer.appendChild(fieldGroup);
    });
    
    // Инициализируем flatpickr для всех полей даты
    document.querySelectorAll('.datepicker').forEach(el => {
        flatpickr(el, {
            dateFormat: "d.m.Y",
            locale: "ru"
        });
    });
}

// Обработчик отправки формы
dynamicForm.addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = {
        university: universitySelect.value,
        supervisor: supervisorSelect.value,
        fields: {}
    };
    
    // Сбор данных полей
    universitiesData[formData.university].fields.forEach(field => {
        formData.fields[field.name] = document.getElementById(field.name).value;
    });
    
    try {
        const token = localStorage.getItem('authToken');
        const headers = {
            'Content-Type': 'application/json',
        };
        
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        const response = await fetch('/forms/api/post/form', {
            method: 'POST',
            headers: headers,
            body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
            throw new Error('Ошибка отправки формы');
        }
        
        const result = await response.json();
        if (result.success) {
            alert('Заявление успешно отправлено!');
            window.location.href = 'success.html';
        } else {
            alert('Ошибка: ' + (result.message || 'Неизвестная ошибка'));
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка при отправке формы: ' + error.message);
    }
});

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    // Инициализация flatpickr для отдельных полей даты (если они есть в HTML)
    flatpickr("#start-date", {
        dateFormat: "d.m.Y",
        onChange: function(selectedDates, dateStr, instance) {
            endDatePicker.set("minDate", dateStr);
        }
    });

    const endDatePicker = flatpickr("#end-date", {
        dateFormat: "d.m.Y"
    });
    
    // Загружаем данные с сервера
    loadUniversities();
    populateSupervisorSelect();
});