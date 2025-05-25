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
        console.error('Ошибка при загрузке университетов:', error);
        alert('Не удалось загрузить список университетов');
    }
}

// Заполнение списка университетов
function populateUniversitySelect() {
    universitySelect.innerHTML = '';

    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'Выберите университет';
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

        const input = document.createElement('input');
        input.id = field.name;
        input.name = field.name;
        input.type = field.type || 'text';
        input.required = !!field.required;
        input.className = 'form-control';

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

// Отправка формы
dynamicForm.addEventListener('submit', async function (e) {
    e.preventDefault();

    const selectedUniversity = universitySelect.value;
    const university = universitiesData[selectedUniversity];

    const infoArray = [];

    university.fields.forEach(field => {
        const value = document.getElementById(field.name).value;
        infoArray.push(value); // просто массив значений
    });

    const formData = {
        institution_id: university.id,
        info: infoArray
    };


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

        if (!response.ok) throw new Error('Ошибка отправки формы');

        const result = await response.json();
        if (response.ok) {
            alert('Заявление успешно отправлено!');
        } else {
            const errorText = await response.text();
            alert('Ошибка: ' + (errorText || 'Неизвестная ошибка'));
        }

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
