// Функция для показа модального окна
function showCreateRoomForm() {
    var modal = document.getElementById("createRoomModal");
    modal.style.display = "block";
}

function hideCreateRoomForm() {
    var modal = document.getElementById("createRoomModal");
    modal.style.display = "none";
}
var autoPlayCheckbox = document.getElementById("autoPlay");

function toggleTimeControl() {
    var timeControl = document.getElementById("timeControl");
    if (autoPlayCheckbox.checked) {
        timeControl.style.display = "block";
    } else {
        timeControl.style.display = "none";
    }
}

// Добавляем слушатель события на изменение чекбокса при загрузке страницы
window.addEventListener('load', function () {
    toggleTimeControl(); // Вызываем функцию, чтобы установить начальное состояние

    // Добавляем слушатель события на изменение чекбокса
    autoPlayCheckbox.addEventListener("change", toggleTimeControl);
});
// Функция для обработки создания комнаты
function createRoom() {
    // Получите данные из формы создания комнаты
    const roomName = document.getElementById('roomName').value;
    const isAuto = document.getElementById('autoPlay').checked;
    const timeOptions = document.getElementById('timeOptions');  // Moved this line outside the if statement
    const time = timeOptions.options[timeOptions.selectedIndex].value;
    const maxPlayers = document.getElementById('maxPlayers').value;

    // Отправить данные на сервер для создания комнаты
    fetch('/api/rooms/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
            roomName: roomName,
            isAuto: isAuto,
            time: time,
            maxPlayers: maxPlayers,
            createRoom: '1', // параметр для различия запроса на создание комнаты
        }),
    })
        .then(response => response.json())
        .then(data => {
            console.log(data);
            if (data.success) {
                // Успешно создано, добавьте комнату к таблице без перезагрузки
                addRoomToTable(data.roomId, roomName, isAuto, time, maxPlayers);
                const err = document.getElementById('error');
                err.innerHTML = "";
                hideCreateRoomForm(); // Скрыть модальное окно
            } else {
                const err = document.getElementById('error');
                err.innerHTML = "";
                console.log("error", data.errors);
                for (let KeyErr in data.errors) {
                    err.innerHTML += data.errors[KeyErr] + '<br>';
                    console.log(KeyErr);
                }
            }
        })
        .catch(error => console.error('Ошибка при создании комнаты: ' + error));
}

// Функция для добавления комнаты к таблице без перезагрузки
function addRoomToTable(id, name, isAuto, time, maxPlayers, organizer) {
    var roomTable = document.getElementById('roomTable');
    var row = roomTable.insertRow(-1);

    var cell1 = row.insertCell(0);
    var cell2 = row.insertCell(1);
    var cell3 = row.insertCell(2);
    var cell4 = row.insertCell(3);
    var cell5 = row.insertCell(4);

    // Create a hyperlink element
    var roomLink = document.createElement('a');
    roomLink.textContent = name;
    roomLink.href = '/rooms/' + name; // Set the href attribute to the room URL
    cell1.appendChild(roomLink);

    cell2.innerHTML = (isAuto ? 'Да' : 'Нет');
    cell3.innerHTML = time;
    cell4.innerHTML = maxPlayers;

    // Создаем кнопку "Удалить" и устанавливаем для нее обработчик события
    var deleteButton = document.createElement('button');
    deleteButton.textContent = 'Удалить';
    deleteButton.addEventListener('click', function () {
        deleteRoom(id); // Вызываем функцию удаления при нажатии на кнопку "Удалить"
    });

    // Добавляем кнопку "Удалить" в ячейку
    cell5.appendChild(deleteButton);
}

// Функция для обновления данных о комнатах на странице
function updateRoomList() {
    fetch('/api/rooms') // Отправить запрос на получение данных
        .then(response => response.json())
        .then(data => {
            var roomTable = document.getElementById('roomTable');

            // Удалить все строки, кроме заголовка
            while (roomTable.rows.length > 1) {
                roomTable.deleteRow(1);
            }

            data.rooms.forEach(room => {
                var row = roomTable.insertRow(-1);
                var cell1 = row.insertCell(0);
                var cell2 = row.insertCell(1);
                var cell3 = row.insertCell(2);
                var cell4 = row.insertCell(3);
                var cell5 = row.insertCell(4);

                // Create a hyperlink element
                var roomLink = document.createElement('a');
                roomLink.textContent = room.Name;
                roomLink.href = '/rooms/' + room.Name; // Set the href attribute to the room URL
                cell1.appendChild(roomLink);

                cell2.innerHTML = (room.Time > 0 ? 'Да' : 'Нет');
                cell3.innerHTML = room.Time;
                cell4.innerHTML = room.Max_partic; // Используем правильное название переменной

                // Создаем кнопку "Удалить" и устанавливаем для нее обработчик события
                var deleteButton = document.createElement('button');
                deleteButton.textContent = 'Удалить';
                deleteButton.addEventListener('click', function () {
                    deleteRoom(room.Name); // Вызываем функцию удаления при нажатии на кнопку "Удалить"
                });

                // Добавляем кнопку "Удалить" в ячейку
                cell5.appendChild(deleteButton);
            });
        })
        .catch(error => console.error('Ошибка при обновлении данных о комнатах: ', error));
}

// Показать/скрыть поле "Время на ход" в зависимости от состояния "Авто-игра"
var autoPlayCheckbox = document.getElementById("autoPlay");
var timeControl = document.getElementById("timeControl");

autoPlayCheckbox.addEventListener("change", function () {
    if (autoPlayCheckbox.checked) {
        timeControl.style.display = "none";
    } else {
        timeControl.style.display = "block";
    }
});

// Вызываем функцию для обновления данных о комнатах при загрузке страницы
window.addEventListener('load', updateRoomList);

// Функция для периодического обновления данных о комнатах
function autoUpdateRoomList() {
    console.log('update')
    updateRoomList(); // Вызываем функцию для обновления данных
    setTimeout(autoUpdateRoomList, 5000); // Вызываем эту же функцию каждые 5 секунд
}

// Вызываем функцию для первоначального обновления данных и запускаем автоматическое обновление
window.addEventListener('load', function () {
    updateRoomList(); // Обновление при загрузке страницы
    autoUpdateRoomList(); // Автоматическое обновление каждые 5 секунд
});
function deleteRoom(roomName) {
    if (confirm('Вы уверены, что хотите удалить эту комнату?')) {
        // Отправить запрос на сервер для удаления комнаты
        fetch('/api/rooms', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                roomName: roomName,
            }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {

                    // Успешно удалено, обновите список комнат
                    updateRoomList();
                } else {
                    console.error('Ошибка при удалении комнаты: ' + data.error);
                }
            })
            .catch(error => console.error('Ошибка при удалении комнаты: ' + error));
    }
}