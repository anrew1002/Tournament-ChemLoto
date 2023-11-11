<?php
// update_room_list.php

// Выполнить подключение к базе данных и запрос для получения данных о комнатах (замените значения на свои)
$dbhost = 'localhost'; // Хост базы данных (в вашем случае, localhost)
            $dbname = 'vloodek20'; // Имя вашей базы данных
            $dbuser = 'root'; // Пользователь базы данных
            $dbpass = ''; // Пароль (пустой, так как у вас нет пароля)

            try {
                $pdo = new PDO("mysql:host=$dbhost;dbname=$dbname", $dbuser, $dbpass);
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

    // Выполните SQL-запрос для выборки данных о комнатах из таблицы rooms
    $query = "SELECT * FROM rooms";
    $stmt = $pdo->query($query);

    // Здесь вы можете форматировать данные и отправлять их на room_list.php
    // Например, создать JSON-ответ и вернуть его клиенту
    $rooms = $stmt->fetchAll(PDO::FETCH_ASSOC);
    echo json_encode($rooms);
} catch (PDOException $e) {
    echo 'Ошибка при получении данных о комнатах: ' . $e->getMessage();
}
?>
