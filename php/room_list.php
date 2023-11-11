<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f5f5f5;
        }
        .close {
        position: absolute;
        top: 10px;
        right: 10px;
        font-size: 108px;
        color: red;
        cursor: pointer;
}
        #header {
            text-align: right;
            padding: 10px;
        }

        #content {
            max-width: 800px;
            margin: 0 auto;
            background-color: #fff;
            border: 1px solid #ccc;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }

        #content h2 {
            text-align: center;
        }

        #roomTable {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }

        #roomTable th, #roomTable td {
            border: 1px solid #ccc;
            padding: 10px;
            text-align: center;
        }

        #roomTable th {
            background-color: #007BFF;
            color: #fff;
        }
        .modal {
            display: none;
            position: fixed;
            z-index: 1;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0,0,0,0.4);
        }

        .modal-content {
            background-color: #fff;
            margin: 15% auto;
            padding: 20px;
            width: 60%;
            border: 1px solid #888;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <?php
    if (session_status() == PHP_SESSION_NONE) {
        session_start();
    }
    // Проверка, если пользователь не вошел в систему, перенаправляем его на страницу входа
    if (!isset($_SESSION['username'])) {
        header('Location: index.php'); 
        exit();
    }
    ?>

    <div id="header">
        <?php
        if (session_status() == PHP_SESSION_NONE) {
            session_start();
        }
        // Показать имя пользователя, если вошли в систему
        if (isset($_SESSION['username'])) {
            echo 'Здравствуйте, ' . $_SESSION['username'] . '!';
        }
        ?>
        <form method="post" action="logout.php" style="display: inline;">
            <button type="submit" name="logout"  style="background-color: red; color: white;">Выход</button>
        </form>
        <form method="post" action="admin_panel.php" style="display: inline;">
            <button type="submit" name="admin_panel">Админ-панель</button>
        </form>
    </div>
    <button onclick="showCreateRoomForm()">Создать комнату</button>

<!-- Модальное окно для создания комнаты -->
<div id="createRoomModal" class="modal">
    <div class="modal-content">
    <span class="close" onclick="hideCreateRoomForm()">&times;</span>
        <h2>Создать комнату</h2>
        <form>
    <label for="roomName">Название комнаты:</label>
    <input type="text" id="roomName" required>

    <label for="autoPlay">Авто-игра:</label>
    <input type="checkbox" id="autoPlay">

    <!-- Этот блок будет появляться при выборе "Авто-игра" -->
    <div id="timeControl" style="display: none;">
        <label>Время на ход (в секундах):</label>
        <select id="timeOptions">
            <option value="20">20 секунд</option>
            <option value="30">30 секунд</option>
            <option value="60">60 секунд</option>
        </select>
    </div>

    <label for="maxPlayers">Макс. количество игроков:</label>
    <input type="number" id="maxPlayers" min="2">

    <button type="button" onclick="createRoom()">Создать</button>
</form>
    </div>
</div>

<script>
    // Функция для показа модального окна
function showCreateRoomForm() {
    var modal = document.getElementById("createRoomModal");
    modal.style.display = "block";
}
function hideCreateRoomForm() {
    var modal = document.getElementById("createRoomModal");
    modal.style.display = "none";
}


// Функция для обработки создания комнаты
function createRoom() {
        // Получите данные из формы создания комнаты
        const roomName = document.getElementById('roomName').value;
        const isAuto = document.getElementById('autoPlay').checked;
        const timeOptions = document.getElementById('timeOptions');
        const time = timeOptions.options[timeOptions.selectedIndex].value;
        const maxPlayers = document.getElementById('maxPlayers').value;

        // Отправить данные на сервер для создания комнаты
        fetch('create_room.php', {
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
                if (data.success) {
                    // Успешно создано, добавьте комнату к таблице без перезагрузки
                    addRoomToTable(data.roomId, roomName, isAuto, time, maxPlayers, data.organizer);
                    hideCreateRoomForm(); // Скрыть модальное окно
                } else {
                    console.error('Ошибка при создании комнаты: ' + data.error);
                }
            })
            .catch(error => console.error('Ошибка при создании комнаты: ' + error));
    }

    // Функция для добавления комнаты к таблице без перезагрузки
    function addRoomToTable(id, name, isAuto, time, maxPlayers, organizer) {
        var roomTable = document.getElementById('roomTable');
        var row = '<tr><td>' + id + '</td><td><a href="rooms/' + id + '.php">' + name + '</a></td><td>' + (isAuto ? 'Да' : 'Нет') + '</td><td>' + time + '</td><td>' + maxPlayers + '</td><td>';
        roomTable.insertAdjacentHTML('beforeend', row);
    }

// Функция для обновления данных о комнатах на странице
function updateRoomList() {
    fetch('update_room_list.php') // Отправить запрос на получение данных
        .then(response => response.json())
        .then(data => {
            var roomTable = document.getElementById('roomTable');

            // Удалить все строки, кроме заголовка
            while (roomTable.rows.length > 1) {
                roomTable.deleteRow(1);
            }

            data.forEach(room => {
                var row = roomTable.insertRow(-1);
                var cell1 = row.insertCell(0);
                var cell2 = row.insertCell(1);
                var cell3 = row.insertCell(2);
                var cell4 = row.insertCell(3);
                var cell5 = row.insertCell(4);



                cell1.innerHTML = room.id;
                cell2.innerHTML = '<a href="rooms/' + room.id + '.php">' + room.name + '</a>';
                cell3.innerHTML = room.is_auto ? 'Да' : 'Нет';
                cell4.innerHTML = room.time;
                cell5.innerHTML = room.max_players;


            });
        })
        .catch(error => console.error('Ошибка при обновлении данных о комнатах: ', error));
}

// Показать/скрыть поле "Время на ход" в зависимости от состояния "Авто-игра"
var autoPlayCheckbox = document.getElementById("autoPlay");
var timeControl = document.getElementById("timeControl");

autoPlayCheckbox.addEventListener("change", function() {
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
window.addEventListener('load', function() {
    updateRoomList(); // Обновление при загрузке страницы
    autoUpdateRoomList(); // Автоматическое обновление каждые 5 секунд
});
var autoPlayCheckbox = document.getElementById("autoPlay");
    var timeControl = document.getElementById("timeControl");

    autoPlayCheckbox.addEventListener("change", function() {
        if (autoPlayCheckbox.checked) {
            timeControl.style.display = "block"; // Показать поле "Время на ход"
        } else {
            timeControl.style.display = "none"; // Скрыть поле "Время на ход"
        }
    });
    function deleteRoom(roomId) {
    if (confirm('Вы уверены, что хотите удалить эту комнату?')) {
        // Отправить запрос на сервер для удаления комнаты
        fetch('delete_room.php', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({
                roomId: roomId,
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
var row = '<tr><td>' + '<a href="rooms/' + room.id + '.php">' + room.name + '</a>' + '</td><td>' + (room.is_auto ? 'Да' : 'Нет') + '</td><td>' + room.time + '</td><td>' + room.max_players + '</td><td><button onclick="deleteRoom(' + room.id + ')">Удалить</button></td></tr>';
roomTable.insertAdjacentHTML('beforeend', row);
</script>

<div id="content">
    <h2>Список комнат для игры</h2>
    <table id="roomTable">
        <tr>
            <th>Название комнаты</th>
            <th>Авто-игра</th>
            <th>Время на ход</th>
            <th>Макс. игроков</th>

        </tr>
    </table>
</div>
    <?php

?>
</body>
</html>