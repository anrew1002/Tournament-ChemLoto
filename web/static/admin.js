function deleteRooms (roomName) {
    console.log(roomName)
    if (confirm('Вы уверены, что хотите удалить эту комнату?')) {
      // Отправить запрос на сервер для удаления комнаты
      fetch('/api/rooms', {
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
          } else {
            console.error('Ошибка при удалении комнаты: ' + data.error)
          }
        })
        .catch(error => console.error('Ошибка при удалении комнаты: ' + error))
    }
}
document.getElementById("roomdel").addEventListener('click',deleteRooms)

function deleteUsers (roomName) {
    console.log(roomName)
    if (confirm('Вы уверены, что хотите удалить эту комнату?')) {
      // Отправить запрос на сервер для удаления комнаты
      fetch('/api/users', {
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
          } else {
            console.error('Ошибка при удалении комнаты: ' + data.error)
          }
        })
        .catch(error => console.error('Ошибка при удалении комнаты: ' + error))
    }
}
document.getElementById("userdel").addEventListener('click',deleteUsers)