//Выбор даты
flatpickr("#start-date", {
    dateFormat: "d.m.Y",
    onChange: function(selectedDates, dateStr, instance) {
        endDatePicker.set("minDate", dateStr);
    }
});

const endDatePicker = flatpickr("#end-date", {
    dateFormat: "d.m.Y"
});

// Скрипт для создания формы в зависимости от ВУЗа
// Пример данных из БД (в реальности будет fetch-запрос)
const universitiesData = {
    "МГУ": {
        "fields": [
            {"name": "teacher", "label": "Преподаватель", "type": "text", "required": true},
            {"name": "practice_start", "label": "Начало практики", "type": "date", "required": true},
            {"name": "practice_end", "label": "Конец практики", "type": "date", "required": true}
        ]
    },
    "ВШЭ": {
        "fields": [
            {"name": "student_name", "label": "ФИО студента", "type": "text", "required": true},
            {"name": "discipline", "label": "Дисциплина", "type": "text", "required": false}
        ]
    }
};

const supervisorData = {
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

// Заполняем список университетов
Object.keys(universitiesData).forEach(universityName => {
    const option = document.createElement('option');
    option.value = universityName;
    option.textContent = universityName;
    universitySelect.appendChild(option);
});

// Заполняем список руководителей практик
Object.values(supervisorData).forEach(supervisorName => {
    const option = document.createElement('option');
    option.value = supervisorName;
    option.textContent = supervisorName;
    supervisorSelect.appendChild(option);
});

// Обработчик изменения выбора университета
universitySelect.addEventListener('change', function() {
    const selectedUniversity = this.value;
    
    // Очищаем предыдущие поля
    formFieldsContainer.innerHTML = '';
    submitBtn.classList.add('hidden');
    
    if (selectedUniversity && universitiesData[selectedUniversity]) {
        // Создаем поля формы согласно конфигурации
        createFormFields(universitiesData[selectedUniversity].fields);
        submitBtn.classList.remove('hidden');
    }
});

// Функция создания полей формы
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
        }
        
        input.id = field.name;
        input.name = field.name;
        input.required = field.required;
        input.className = 'form-control';
        
        fieldGroup.appendChild(label);
        fieldGroup.appendChild(input);
        formFieldsContainer.appendChild(fieldGroup);
    });
}

// Обработчик отправки формы
dynamicForm.addEventListener('submit', function(e) {
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
            const response = await fetch('/api/submit-student-form', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });
            
            const result = await response.json();
            if (result.success) {
                alert('Заявление успешно отправлено!');
                window.location.href = 'success.html';
            } else {
                alert('Ошибка: ' + (result.message || 'Неизвестная ошибка'));
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Ошибка при отправке формы');
        }
    });
});

document.addEventListener('DOMContentLoaded', function() {
    // Загрузка данных форм (из localStorage или с сервера)
    loadUniversityForms();
    
    // Заполнение списка университетов
    universitySelect();
});

// Загрузка из localStorage
function loadUniversityForms() {
    const savedForms = localStorage.getItem('universityForms');
    if (savedForms) {
        universitiesData = JSON.parse(savedForms);
    }
}