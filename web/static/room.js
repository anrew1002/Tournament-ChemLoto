document.addEventListener('DOMContentLoaded', function () {
  
  var isAdmin = document.getElementById('isAdmin').textContent
  console.log(isAdmin)

  function fetchAndUpdateTopPlayers() {
    fetch('/api/users')
      .then(response => response.json())
      .then(data => {
        const topPlayersList = document.getElementById('topPlayersList')
        const roomName = document.getElementById('room-title')

        // Clear existing list
        topPlayersList.innerHTML = ''

        // Фильтруем пользователей по комнате (замените "Название Комнаты" на фактическое название комнаты)
        const usersInRoom = data.users.filter(
          user => user.Room === roomName.textContent
        )

        // Sort users by score in descending order
        usersInRoom.sort((a, b) => b.Score - a.Score)

        // Loop through the sorted user data and update the table
        usersInRoom.forEach(user => {
          const listItem = document.createElement('li')
          listItem.textContent = `${user.Username} - ${user.Score} очков`

          // Add a click event listener to each player name
          listItem.addEventListener('click', function () {
            if (isAdmin === 'true') {
              console.log('isAdmin')
              // Сохраните имя пользователя в глобальной переменной
              selectedUsername = user.Username
              openModal(user) // Open the modal for the selected player
            }
          })

          topPlayersList.appendChild(listItem)
        })
      })
      .catch(error => console.error('Error fetching user data:', error))
  }

  // Вызовите функцию без передачи значения isAdmin
  fetchAndUpdateTopPlayers()

  // Set up an interval to periodically update the top players table
  setInterval(fetchAndUpdateTopPlayers, 1000) // Update every minute (adjust as needed)
  var selectedUsername
  function openModal(player) {
    console.log('da')
    var modal = document.getElementById('myModal')
    var modalTitle = document.getElementById('modalTitle')

    // Изменяем текст заголовка в соответствии с именем игрока
    modalTitle.textContent = 'Начислить очки игроку: ' + player.Username

    modal.style.display = 'block'
    document.getElementById('playerName').textContent =
      'Игрок: ' + player.Username

    // Сохраните имя пользователя в глобальной переменной
    selectedUsername = player.Username
    const encodedUsername = encodeURIComponent(selectedUsername)
  }

  function closeModal() {
    var modal = document.getElementById('myModal')
    modal.style.display = 'none'
  }

  // Обработчик клика по кнопке закрытия (крестик)
  document.querySelector('.close-btn').addEventListener('click', closeModal)

  function getSelectedScore(block) {
    // Get the value of the selected radio button in a block
    var selectedRadio = block.querySelector('input:checked')
    return selectedRadio ? selectedRadio.value : null
  }

  function sendScores() {
    // Получите значения из полей формы
    const alphaScore = getSelectedScore(document.getElementById('alphaBlock'))
    const betaScore = getSelectedScore(document.getElementById('betaBlock'))
    const gammaScore = getSelectedScore(document.getElementById('gammaBlock'))
    const penaltyScore = getSelectedScore(document.getElementById('penaltyBlock'))
    const manualScore = document.getElementById('manualScore').value || 0; // Default to 0 if nothing is entered manually

    // Замените 'selectedUsername' на актуальное значение имени пользователя
    const encodedUsername = encodeURIComponent(selectedUsername);

    // Создаем объект сообщения для каждого типа score
    const messages = [];

    if (alphaScore !== 0) {
      messages.push({
        type: 'score_up',
        struct: {
          field: "alpha",
          target: selectedUsername,
          score: parseInt(alphaScore)
        }
      });
    }

    if (betaScore !== 0) {
      messages.push({
        type: 'score_up',
        struct: {
          field: "beta",
          target: selectedUsername,
          score: parseInt(betaScore)
        }
      });
    }

    if (gammaScore !== 0) {
      messages.push({
        type: 'score_up',
        struct: {
          field: "gamma",
          target: selectedUsername,
          score: parseInt(gammaScore)
        }
      });
    }

    if (penaltyScore !== 0) {
      messages.push({
        type: 'score_up',
        struct: {
          field: "penalty",
          target: selectedUsername,
          score: parseInt(penaltyScore)
        }
      });
    }

    if (manualScore !== 0) {
      messages.push({
        type: 'score_up',
        struct: {
          field: "manual",
          target: selectedUsername,
          score: parseInt(manualScore)
        }
      });
    }

    // Отправляем сообщения на сервер
    messages.forEach(message => {
      socket.send(JSON.stringify(message));
    });

    // // Отправьте данные на сервер
    // fetch('/api/users/' + encodedUsername, {
    //   method: 'POST',
    //   headers: {
    //     'Content-Type': 'application/x-www-form-urlencoded'
    //   },
    //   body: new URLSearchParams({
    //     score: alphaScore
    //   })
    // })
    //   .then(response => response.json())
    //   .then(data => {
    //     console.log(data)
    //     if (data.success) {
    //       // Обработайте успешный ответ от сервера, если необходимо
    //       console.log('Очки успешно отправлены')
    //     } else {
    //       // Обработайте ошибку от сервера, если необходимо
    //       console.error('Ошибка при отправке очков:', data.errors)
    //     }
    //   })
    //   .catch(error => console.error('Ошибка при отправке очков: ' + error))
  }

  document.getElementById('modalButton').addEventListener('click', function () {
    // Call the sendScores function
    sendScores()
    // Close the modal or perform other actions as needed
    closeModal()
  })
})
document.getElementById('chatToggle').onclick = function () {
  var chatElement = document.querySelector('.chat')

  if (chatElement.style.display === 'none') chatElement.style.display = 'flex'
  else chatElement.style.display = 'none'
}
document.getElementById('topToggle').onclick = function () {
  var topElement = document.querySelector('.top-players')

  if (topElement.style.display === 'none') topElement.style.display = 'block'
  else topElement.style.display = 'none'
}
// function raiseHand() {

//     // Создаем один раз экземпляр WebSocket
//     // const wsUrl = 'ws://127.0.0.1/ws';
//     // const socket = new WebSocket(wsUrl);

//     // // Слушаем событие открытия соединения
//     // socket.addEventListener('open', (event) => {
//     //     console.log('WebSocket connection opened:', event);
//     // });

//     // // Слушаем событие закрытия соединения
//     // socket.addEventListener('close', (event) => {
//     //     console.log('WebSocket connection closed:', event);
//     // });

//     // // Слушаем событие ошибки
//     // socket.addEventListener('error', (error) => {
//     //     console.error('WebSocket error:', error);
//     // });

//     // Слушаем события от сервера
//     socket.addEventListener('message', (event) => {
//         const data = JSON.parse(event.data);

//         // Проверяем тип сообщения
//         if (data.type === 'raiseHandNotification') {
//             // Ваш код для обработки уведомления о поднятии руки
//             console.log('Кто-то поднял руку!');
//             // Здесь вы можете выполнить какие-то действия для отображения уведомления на странице
//         }
//     });

//     // Ваша функция для отправки уведомления о поднятии руки

// }
function raiseHand() {
  // Создаем объект сообщения
  const message = {
    type: 'raise_hand'
  }


  // Отправляем сообщение на сервер
  socket.send(JSON.stringify(message))
}

function getElement() {
  // Создаем объект сообщения
  const message = {
    type: 'get_element'
  }
  // Отправляем сообщение на сервер
  socket.send(JSON.stringify(message))
}
function startGame() {
  // Ваш код для начала игры
  const message = {
    type: 'start_game'
  }
  
  // Отправляем сообщение на сервер
  socket.send(JSON.stringify(message))

}

console.log(timer)
function stopGame() {
  // Ваш код для начала игры
  const message = {
    type: 'raise_hand'
  }
  var startButton = document.getElementById('continueButton');
  startButton.style.display = 'block';

  var stopButton = document.getElementById('stopButton');
  stopButton.style.display = 'none';
  // Отправляем сообщение на сервер
  socket.send(JSON.stringify(message))

}
function continueGame() {
  // Ваш код для начала игры
  const message = {
    type: 'continue_game'
  }
  var stopButton = document.getElementById('stopButton');
  stopButton.style.display = 'block';

  var startButton = document.getElementById('continueButton');
  startButton.style.display = 'none';
  // Отправляем сообщение на сервер
  socket.send(JSON.stringify(message))

}