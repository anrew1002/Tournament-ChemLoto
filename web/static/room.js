document.addEventListener("DOMContentLoaded", function () {
    var timerElement = document.querySelector(".timer");
    var imageElement = document.getElementById("elementImage");
    var initialTime = 20;

    function updateTimer() {
        timerElement.textContent = formatTime(initialTime);

        // Добавлено условие для мерцания
        if (initialTime <= 5 && initialTime % 2 === 0) {
            imageElement.classList.add("flash");
        } else {
            imageElement.classList.remove("flash");
        }

        // Добавлено условие для анимации смены картинки
        if (initialTime <= 0) {
            imageElement.src = "chem_el2.png";
            resetTimer();
        }

        initialTime--;
        setTimeout(updateTimer, 1000);
    }

    function formatTime(seconds) {
        var minutes = Math.floor(seconds / 60);
        var remainingSeconds = seconds % 60;

        var formattedTime =
            pad(minutes, 2) + ":" + pad(remainingSeconds, 2);

        return formattedTime;
    }

    function pad(number, length) {
        var str = String(number);
        while (str.length < length) {
            str = "0" + str;
        }
        return str;
    }

    function resetTimer() {
        initialTime = 20;
    }

    updateTimer();
    function fetchAndUpdateTopPlayers() {
        fetch("/api/users")
            .then(response => response.json())
            .then(data => {
                const topPlayersList = document.getElementById("topPlayersList");

                // Clear existing list
                topPlayersList.innerHTML = "";

                // Фильтруем пользователей по комнате (замените "Название Комнаты" на фактическое название комнаты)
                const usersInRoom = data.users.filter(user => user.Room === "");

                // Sort users by score in descending order
                usersInRoom.sort((a, b) => b.Score - a.Score);

                // Loop through the sorted user data and update the table
                usersInRoom.forEach(user => {
                    const listItem = document.createElement("li");
                    listItem.textContent = `${user.Username} - ${user.Score} очков`;
                    topPlayersList.appendChild(listItem);
                });
            })
            .catch(error => console.error("Error fetching user data:", error));
    }


    // ... (existing code)

    // Fetch and update top players initially
    fetchAndUpdateTopPlayers();

    // Set up an interval to periodically update the top players table
    setInterval(fetchAndUpdateTopPlayers, 1000); // Update every minute (adjust as needed)
});
document.getElementById('chatToggle').onclick = function () {
    var chatElement = document.querySelector('.chat');

    if (chatElement.style.display === 'none') chatElement.style.display = 'flex';
    else chatElement.style.display = 'none';
};
document.getElementById('topToggle').onclick = function () {
    var topElement = document.querySelector('.top-players');

    if (topElement.style.display === 'none') topElement.style.display = 'block';
    else topElement.style.display = 'none';
};
function raiseHand() {
    // Добавьте следующий код для отправки уведомления на сервер WebSocket
    const message = {
        type: 'raiseHand',
    }

    // Создаем один раз экземпляр WebSocket
    const wsUrl = 'ws://127.0.0.1/ws';
    const socket = new WebSocket(wsUrl);

    // Слушаем событие открытия соединения
    socket.addEventListener('open', (event) => {
        console.log('WebSocket connection opened:', event);
    });

    // Слушаем событие закрытия соединения
    socket.addEventListener('close', (event) => {
        console.log('WebSocket connection closed:', event);
    });

    // Слушаем событие ошибки
    socket.addEventListener('error', (error) => {
        console.error('WebSocket error:', error);
    });

    // Слушаем события от сервера
    socket.addEventListener('message', (event) => {
        const data = JSON.parse(event.data);

        // Проверяем тип сообщения
        if (data.type === 'raiseHandNotification') {
            // Ваш код для обработки уведомления о поднятии руки
            console.log('Кто-то поднял руку!');
            // Здесь вы можете выполнить какие-то действия для отображения уведомления на странице
        }
    });

    // Ваша функция для отправки уведомления о поднятии руки
    function raiseHand() {
        // Создаем объект сообщения
        const message = {
            type: 'raiseHand',
        };

        // Отправляем сообщение на сервер
        socket.send(JSON.stringify(message));
    }
}