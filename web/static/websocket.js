'use strict'
var socket

document.addEventListener('DOMContentLoaded', e => {
  // get important elements
  var isAdmin = document.getElementById('isAdmin').textContent
  const dialod_window = document.querySelector('.messages')
  const enc = new TextDecoder('utf-8')
  const status = document.getElementById('status')
  const message_template = document
    .getElementById('start_message')
    .cloneNode(true)
  let accounts_map = new Map()
  setWebsocket()

  function setWebsocket() {
    socket = new WebSocket('ws://127.0.0.1:80/ws')
    socket.onmessage = function (event) {
      messageHandler(JSON.parse(event.data))
    }
    socket.onclose = function (e) {
      console.log(
        'Socket is closed. Reconnect will be attempted in 1 second.',
        e.reason
      )
      // status.textContent = 'Офлайн'
      setTimeout(function () {
        setWebsocket()
      }, 1000)
    }
    socket.onopen = function (e) {
      // status.textContent = 'Онлайн'
    }

    return socket
  }
  document.forms['publish'].onsubmit = function (e) {
    e.preventDefault()
    // console.log(this.message.value)
    const chat_message = {
      type: 'chat_text',
      // struct: JSON.stringify(this.message.value)
      struct: this.message.value
    }
    let outgoingMessage = this.message.value
    this.message.value = ''
    if (outgoingMessage != '') {
      console.log(chat_message)
      socket.send(JSON.stringify(chat_message))
    }
    // return false
  }
  var timer = null;
  var intervalId = null;
  var isTimerPaused = false;
  let isStopButtonHidden = false;
  async function messageHandler(data) {
    switch (data.type) {
      case 'chat_text':
        textMessageHandler(data);
        break;
      case 'raise_hand':
        raiseHandNotification(data.struct.sender);
        pauseTimer();
        if (!isStopButtonHidden && isAdmin ==='true') {
        var startButton = document.getElementById('continueButton');
        startButton.style.display = 'block';

        var stopButton = document.getElementById('stopButton');
        stopButton.style.display = 'none';
        isStopButtonHidden = true;
        }
        break;
      case 'send_element':
        handleElementResponse(data.struct.element, data.struct.last_elements);
        // Сразу сбрасываем и запускаем таймер
        if (timer == 0) {

        } else {
          resetAndStartTimer(timer);
        }
        break;
      case 'start_game':
        startGameHandler();
        resumeTimer()
        if (timer == 0 && isAdmin === 'true') {
          var stopButton = document.getElementById('stopButton');
          stopButton.style.display = 'none';
        } else {
          timerHandler(data.struct.Time);
        }
        timer = data.struct.Time;
        break;
      case 'init_connection':
        timer = data.struct.Time;
        if (data.struct.Started == true) {
          startGameHandler();
          console.log(timer)
          if (timer === 0 && !isStopButtonHidden && isAdmin === 'true') {
            var stopButton = document.getElementById('stopButton');
            stopButton.style.display = 'none';
            
            // Устанавливаем флаг в true после первого скрытия
            isStopButtonHidden = true;
        } else if (timer !== 0) {
            timerHandler(data.struct.Time);
        }
          if (data.struct.Paused === true) {
            if (isAdmin === 'true'){
              var startButton = document.getElementById('continueButton');
            startButton.style.display = 'block';

            var stopButton = document.getElementById('stopButton');
            stopButton.style.display = 'none';
            }
            
            pauseTimer()
          }
          console.log(data.struct.last_elements)
          handleElementResponse(data.struct.last_elements[data.struct.last_elements.length - 1], data.struct.last_elements);
        }
        break;
      default:
        console.log('Неизвестный тип сообщения от сервера ', data.type);
        break;
    }
  }

  function timerHandler(time) {
    if (time == 0) {
      return;
    }

    var timerElement = document.querySelector('.timer');
    var imageElement = document.getElementById('elementImage');
    var initialTime = time;

    function updateTimer() {
      if (isTimerPaused) {
        // Если таймер приостановлен, не обновляем его
        return;
      }

      timerElement.textContent = formatTime(initialTime);

      if (initialTime <= 5 && initialTime % 2 === 0) {
        imageElement.classList.add('flash');
      } else {
        imageElement.classList.remove('flash');
      }

      if (initialTime <= 0) {
        resetAndStartTimer(time);
      }

      initialTime--;
    }

    function formatTime(seconds) {
      var minutes = Math.floor(seconds / 60);
      var remainingSeconds = seconds % 60;
      var formattedTime = pad(minutes, 2) + ':' + pad(remainingSeconds, 2);
      return formattedTime;
    }

    function pad(number, length) {
      var str = String(number);
      while (str.length < length) {
        str = '0' + str;
      }
      return str;
    }

    updateTimer();
    // Очищаем предыдущий interval перед установкой нового
    clearInterval(intervalId);
    // Устанавливаем новый interval
    intervalId = setInterval(updateTimer, 1000);
  }

  function resetAndStartTimer(time) {
    // Здесь должно быть назначение начального времени, если это необходимо.
    timer = time;
    // Сбросить таймер и запустить его заново
    timerHandler(timer);
  }
  function pauseTimer() {
    isTimerPaused = true;
  }
  function resumeTimer() {
    isTimerPaused = false;
  }




  function startGameHandler() {
    console.log('startGameHandler called');
    var isAdmin = document.getElementById('isAdmin').textContent
    if (timer !== 0 && isAdmin == 'true') {

      var stopButton = document.getElementById('stopButton');
      stopButton.style.display = 'block';

    }

    var element = document.getElementById('elementImage');
    // Изменяем свойство display на 'block'
    element.style.display = 'block';
    // Покажите кнопку "Вытащить новый элемент"
    if (timer === 0) {

      var getElementButton = document.querySelector('.admin-btn[onclick="getElement()"]');
      if (getElementButton) {
        console.log('Showing getElementButton');
        getElementButton.style.display = 'block';
      } else {
        console.log('getElementButton not found');
      }
    }
    // Скрыть кнопку "Начать игру"
    var startGameButton = document.querySelector('.admin-btn.start-game-btn');
    if (startGameButton) {
      console.log('Hiding startGameButton');
      startGameButton.style.display = 'none';
    } else {
      console.log('startGameButton not found');
    }

    // Другие действия, если необходимо
  }





  let currentElementIndex = 5; // Variable to store the index of the current element, starting from the last cell
  let currentElement = ''; // Variable to store the current element

  // Assume this function is called when you receive the element data
  function handleElementResponse(element, lastElements) {
    const elementImage = document.getElementById('elementImage');
    // Show the last-elements container if it's not already visible
    const lastElementsContainer = document.getElementById('lastElementsContainer');
    if (lastElementsContainer.style.display === 'none') {
      lastElementsContainer.style.display = 'block';
    }

    // Update the elementImage source based on the received element
    elementImage.src = `../items/${element}.svg`;

    // Update the last element images dynamically
    for (let i = 0; i < lastElements.length; i++) {
      const lastElementImage = document.getElementById(`element${i + 1}`);
      lastElementImage.src = `../items/${lastElements[i] || 'empty'}.svg`;
    }
  }

  function updateLastElementImages() {
    for (let i = 5; i > 1; i--) {
      const currentElementImage = document.getElementById(`element${i}`);
      const previousElementImage = document.getElementById(`element${i - 1}`);
      currentElementImage.src = previousElementImage.src;
    }

    // Check if the first element matches the current element or is empty
    const firstElementImage = document.getElementById('element1');
    if (firstElementImage.src !== `../items/${currentElement}.svg` && currentElement !== '') {
      firstElementImage.src = `../items/${currentElement}.svg`;
    }
  }





  function raiseHandNotification(username) {
    const notificationContainer = document.getElementById(
      'notification-container'
    )
    const notificationText = document.getElementById('notification-text')

    // Устанавливаем текст уведомления с именем пользователя
    notificationText.textContent = `${username} поднял руку!`

    // Показываем уведомление
    notificationContainer.classList.add('show')

    // Проигрываем звук
    playNotificationSound()

    // Скрываем уведомление через 5 секунд
    setTimeout(function () {
      notificationContainer.classList.remove('show')
    }, 8000)
  }
  function playNotificationSound() {
    var notificationSound = document.getElementById('notification-sound')
    notificationSound.play()
  }
  function textMessageHandler(message) {
    console.log(message)
    let messageElem = message_template.cloneNode(true)
    // message.struct.payload = enc.decode(
    //   base64ToArrayBuffer(message.struct.payload)
    // )

    messageElem.querySelector('.message__author').textContent =
      message.struct.sender
    messageElem.querySelector('.message__field').textContent =
      message.struct.payload
    dialod_window.append(messageElem)
  }
  async function fetchUsers(accounts_map) {
    try {
      let response = await fetch('/all_users')
      let accounts = await response.json()
      //TODO: нужно сделать умнее
      accounts.forEach(user => {
        accounts_map.set(user.ID, user)
      })
      return accounts_map
    } catch (err) {
      // перехватит любую ошибку в блоке try: и в fetch, и в response.json
      alert(err)
    }
  }

  function base64ToArrayBuffer(base64) {
    var binaryString = atob(base64)
    var bytes = new Uint8Array(binaryString.length)
    for (var i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }
    return bytes.buffer
  }
})