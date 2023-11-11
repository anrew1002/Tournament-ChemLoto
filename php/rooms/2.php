<div>
    <h1>Добро пожаловать в комнату 123</h1>
    <p>Максимальное количество игроков: 2</p>
    <div id="numberValue">0</div>
    <!-- Дополнительное содержимое вашей комнаты -->
</div>
<script>
const eventSource = new EventSource('../live_update.php');

eventSource.onmessage = function (event) {
    const data = JSON.parse(event.data);
    const randomNumber = data.randomNumber;

    // Обновить значение числа на странице
    const numberValue = document.getElementById('numberValue');
    numberValue.innerText = randomNumber;
};
</script>