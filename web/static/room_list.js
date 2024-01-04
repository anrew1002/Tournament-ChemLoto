// Функция для показа модального окна
function showCreateRoomForm () {
  var modal = document.getElementById('createRoomModal')
  modal.style.display = 'block'
}

function hideCreateRoomForm () {
  var modal = document.getElementById('createRoomModal')
  modal.style.display = 'none'
}
var autoPlayCheckbox = document.getElementById('autoPlay')

function toggleTimeControl () {
  var timeControl = document.getElementById('timeControl')
  if (autoPlayCheckbox.checked) {
    timeControl.style.display = 'block'
  } else {
    timeControl.style.display = 'none'
  }
}

// Добавляем слушатель события на изменение чекбокса при загрузке страницы
window.addEventListener('load', function () {
  toggleTimeControl() // Вызываем функцию, чтобы установить начальное состояние

  // Добавляем слушатель события на изменение чекбокса
  autoPlayCheckbox.addEventListener('change', toggleTimeControl)
})
// Функция для обработки создания комнаты
function createRoom () {
  // Получите данные из формы создания комнаты
  const roomName = document.getElementById('roomName').value
  const isAuto = document.getElementById('autoPlay').checked
  const timeOptions = document.getElementById('timeOptions')
  const time = timeOptions.options[timeOptions.selectedIndex].value
  const maxPlayers = document.getElementById('maxPlayers').value

  // Получите данные о количестве элементов
  const elementCounts = {}
  const chemicalElements = document.querySelectorAll('.chemical-element')
  chemicalElements.forEach(element => {
    const elementId = element.querySelector('.element-img').alt
    const elementCountInput = element.querySelector('input[type="number"]')
    const elementCount = parseInt(elementCountInput.value, 10)
    elementCounts[elementId] = elementCount
  })

  // Отправить данные на сервер для создания комнаты
  fetch('/api/rooms/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json' // Изменено на application/json
    },
    body: JSON.stringify({
      roomName: roomName,
      isAuto: isAuto,
      time: parseInt(time),
      maxPlayers: parseInt(maxPlayers),
      createRoom: '1',
      elementCounts: elementCounts
    })
  })
    .then(response => response.json())
    .then(data => {
      console.log(data)
      if (data.success) {
        // Успешно создано, добавьте комнату к таблице без перезагрузки
        addRoomToTable(data.roomId, roomName, isAuto, time, maxPlayers)
        const err = document.getElementById('error')
        err.innerHTML = ''
        hideCreateRoomForm() // Скрыть модальное окно
      } else {
        const err = document.getElementById('error')
        err.innerHTML = ''
        console.log('error', data.errors)
        for (let KeyErr in data.errors) {
          err.innerHTML += data.errors[KeyErr] + '<br>'
          console.log(KeyErr)
        }
      }
    })
    .catch(error => console.error('Ошибка при создании комнаты: ' + error))
}

// Функция для добавления комнаты к таблице без перезагрузки
function addRoomToTable (id, name, isAuto, time, maxPlayers, organizer) {
  var roomTable = document.getElementById('roomTable')
  var row = roomTable.insertRow(-1)

  var cell1 = row.insertCell(0)
  var cell2 = row.insertCell(1)
  var cell3 = row.insertCell(2)
  var cell4 = row.insertCell(3)
  var cell5 = row.insertCell(4)

  // Create a hyperlink element
  var roomLink = document.createElement('a')
  roomLink.textContent = name
  roomLink.href = '/rooms/' + encodeURI(name) // Set the href attribute to the room URL
  cell1.appendChild(roomLink)

  cell2.innerHTML = isAuto ? 'Да' : 'Нет'
  cell3.innerHTML = time
  cell4.innerHTML = maxPlayers

  // Создаем кнопку "Удалить" и устанавливаем для нее обработчик события
  var deleteButton = document.createElement('button')
  deleteButton.textContent = 'Удалить'
  deleteButton.addEventListener('click', function () {
    console.log(room.name)
    deleteRoom(room.name) // Используйте room.roomName
  })

  // Добавляем кнопку "Удалить" в ячейку
  cell5.appendChild(deleteButton)
}
// Функция для обновления данных о комнатах на странице
function updateRoomList () {
  var delete_th = document.getElementById('if_admin')

  fetch('/api/rooms') // Отправить запрос на получение данных
    .then(response => response.json())
    .then(data => {
      var roomTable = document.getElementById('roomTable')

      // Удалить все строки, кроме заголовка
      while (roomTable.rows.length > 1) {
        roomTable.deleteRow(1)
      }

      data.rooms.forEach(room => {
        var row = roomTable.insertRow(-1)
        var cell1 = row.insertCell(0)
        var cell2 = row.insertCell(1)
        var cell3 = row.insertCell(2)
        var cell4 = row.insertCell(3)
        if (delete_th != null) {
          var cell5 = row.insertCell(4)
        }

        // Create a hyperlink element
        var roomLink = document.createElement('a')
        roomLink.textContent = room.roomName
        roomLink.href = '/rooms/' + encodeURI(room.roomName) // Set the href attribute to the room URL
        cell1.appendChild(roomLink)
        cell2.innerHTML = room.time > 0 ? 'Да' : 'Нет'
        cell3.innerHTML = room.time
        cell4.innerHTML = room.maxPlayers // Используем правильное название переменной

        if (delete_th != null) {
          // Создаем кнопку "Удалить" и устанавливаем для нее обработчик события
          var deleteButton = document.createElement('button')
          deleteButton.textContent = 'Удалить'
          deleteButton.addEventListener('click', function () {
            console.log(room.roomName)
            deleteRoom(room.roomName) // Вызываем функцию удаления при нажатии на кнопку "Удалить"
          })

          // Добавляем кнопку "Удалить" в ячейку
          cell5.appendChild(deleteButton)
        }
      })
    })
    .catch(error =>
      console.error('Ошибка при обновлении данных о комнатах: ', error)
    )
}

// Показать/скрыть поле "Время на ход" в зависимости от состояния "Авто-игра"
var autoPlayCheckbox = document.getElementById('autoPlay')
var timeControl = document.getElementById('timeControl')

autoPlayCheckbox.addEventListener('change', function () {
  if (autoPlayCheckbox.checked) {
    timeControl.style.display = 'none'
  } else {
    timeControl.style.display = 'block'
  }
})

// Вызываем функцию для обновления данных о комнатах при загрузке страницы
window.addEventListener('load', updateRoomList)

// Функция для периодического обновления данных о комнатах
function autoUpdateRoomList () {
  console.log('update')
  updateRoomList() // Вызываем функцию для обновления данных
  setTimeout(autoUpdateRoomList, 5000) // Вызываем эту же функцию каждые 5 секунд
}

// Вызываем функцию для первоначального обновления данных и запускаем автоматическое обновление
window.addEventListener('load', function () {
  updateRoomList() // Обновление при загрузке страницы
  autoUpdateRoomList() // Автоматическое обновление каждые 5 секунд
})

function deleteRoom (roomName) {
  console.log(roomName)
  if (confirm('Вы уверены, что хотите удалить эту комнату?')) {
    // Отправить запрос на сервер для удаления комнаты
    fetch('/api/rooms/' + encodeURI(roomName), {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        roomName: roomName
      })
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          // Успешно удалено, обновите список комнат
          updateRoomList()
        } else {
          console.error('Ошибка при удалении комнаты: ' + data.error)
        }
      })
      .catch(error => console.error('Ошибка при удалении комнаты: ' + error))
  }
}
