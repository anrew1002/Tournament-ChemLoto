'use strict'
var socket
document.addEventListener('DOMContentLoaded', e => {
  // get important elements
  const dialod_window = document.querySelector('.messages')
  const enc = new TextDecoder('utf-8')
  const status = document.getElementById('status')
  const message_template = document
    .getElementById('start_message')
    .cloneNode(true)
  let accounts_map = new Map()
  setWebsocket()

  function setWebsocket () {
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
  document.forms['publish'].onsubmit = function () {
    e.preventDefault()
    // console.log(this.message.value)
    // const chat_message = {
    //   type: 'chat_text',
    //   struct: JSON.stringify(this.message.value)
    // }
    // // let outgoingMessage = this.message.value
    // this.message.value = ''
    // if (outgoingMessage != '') {
    //   // console.log(outgoingMessage)
    //   socket.send(chat_message)
    // }
    // return false
  }

  async function messageHandler (data) {
    // if (!accounts_map.has(data.author_id)) {
    //   accounts_map = await fetchUsers(accounts_map)
    //   console.log('new account data', accounts_map)
    // }
    switch (data.type) {
      case 'chat_text':
        textMessageHandler(data)
        break
      case 'raiseHand':
        console.log(data)
        raiseHandNotification(data.struct.sender) // Передаем имя пользователя

        break
      default:
        console.log('Undefined message type from server')
        break
    }
  }
  function raiseHandNotification (username) {
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
  function playNotificationSound () {
    var notificationSound = document.getElementById('notification-sound')
    notificationSound.play()
  }
  function textMessageHandler (message) {
    console.log(message)
    let messageElem = message_template.cloneNode(true)
    message.struct.payload = enc.decode(
      base64ToArrayBuffer(message.struct.payload)
    )

    messageElem.querySelector('.message__author').textContent =
      message.struct.sender
    messageElem.querySelector('.message__field').textContent =
      message.struct.payload
    dialod_window.append(messageElem)
  }
  async function fetchUsers (accounts_map) {
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

  function base64ToArrayBuffer (base64) {
    var binaryString = atob(base64)
    var bytes = new Uint8Array(binaryString.length)
    for (var i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }
    return bytes.buffer
  }
})
