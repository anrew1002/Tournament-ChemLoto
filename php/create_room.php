<?php
session_start();

if ($_SERVER["REQUEST_METHOD"] == "POST") {
    // Получить данные из POST-запроса
    $roomName = $_POST["roomName"];
    $isAuto = isset($_POST["isAuto"]) ? 1 : 0;
    $time = $_POST["time"];
    $maxPlayers = $_POST["maxPlayers"];

    // Выполнить подключение к базе данных (замените значения на свои)
    $dbhost = 'localhost'; // Хост базы данных (в вашем случае, localhost)
            $dbname = 'vloodek20'; // Имя вашей базы данных
            $dbuser = 'root'; // Пользователь базы данных
            $dbpass = ''; // Пароль (пустой, так как у вас нет пароля)

            try {
                $pdo = new PDO("mysql:host=$dbhost;dbname=$dbname", $dbuser, $dbpass);
        $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

        // Выполнить SQL-запрос для вставки данных в таблицу rooms
        $query = "INSERT INTO rooms (name, is_auto, time, max_players) VALUES (:name, :is_auto, :time, :max_players)";
        $stmt = $pdo->prepare($query);
        $stmt->bindParam(':name', $roomName);
        $stmt->bindParam(':is_auto', $isAuto);
        $stmt->bindParam(':time', $time);
        $stmt->bindParam(':max_players', $maxPlayers);
        $stmt->execute();

        // Получите ID последней вставленной записи
        $roomId = $pdo->lastInsertId();

        // Создайте папку "rooms", если её нет
        $roomsDirectory = 'rooms';
        if (!file_exists($roomsDirectory)) {
            mkdir($roomsDirectory, 0755, true);
        }

        // Создайте файл с именем, соответствующим ID комнаты и вставьте JavaScript-код
        $filePath = $roomsDirectory . '/' . $roomId . '.php';

        $jsCode = <<<EOD
<div>
    <h1>Добро пожаловать в комнату {$roomName}</h1>
    <p>Максимальное количество игроков: {$maxPlayers}</p>
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
EOD;

        file_put_contents($filePath, $jsCode);

        // Верните JSON-ответ с ID комнаты
        $response = [
            'success' => true,
            'roomId' => $roomId,
        ];
        echo json_encode($response);
    } catch (PDOException $e) {
        // Верните JSON-ответ с ошибкой
        $response = [
            'success' => false,
            'error' => 'Ошибка при создании комнаты: ' . $e->getMessage(),
        ];
        echo json_encode($response);
    }
}
?>
