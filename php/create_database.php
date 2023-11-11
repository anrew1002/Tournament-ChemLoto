<?php

session_start();

if (!isset($_SESSION['isOrganizer']) || !$_SESSION['isOrganizer']) {
    die("Доступ запрещен. Пожалуйста, войдите как организатор.");
}
try {
    // Параметры подключения к MySQL серверу
    $dbhost = 'localhost'; // Хост базы данных (в вашем случае, localhost)
            $dbname = 'vloodek20'; // Имя вашей базы данных
            $dbuser = 'root'; // Пользователь базы данных
            $dbpass = ''; // Пароль (пустой, так как у вас нет пароля)

            try {
                $pdo = new PDO("mysql:host=$dbhost;dbname=$dbname", $dbuser, $dbpass);
        $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    } catch (PDOException $e) {

        $message = 'Ошибка подключения: ' . $e->getMessage();
    }

    // Устанавливаем PDO атрибут, чтобы выбрасывать исключения в случае ошибок
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

    // Создание базы данных
    $sql = "CREATE DATABASE IF NOT EXISTS vloodek20";
    $pdo->exec($sql);
    echo "База данных успешно создана.<br>";

     // Удаление старой таблицы "players", если она существует
    $sql = "DROP TABLE IF EXISTS players";
    $pdo->exec($sql);

    // Подключение к созданной базе данных
    $pdo->exec("USE vloodek20");

    // Создание таблицы для игроков
    $sql = "CREATE TABLE IF NOT EXISTS players (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        score INT NOT NULL DEFAULT 0,
        is_organizer BOOLEAN NOT NULL
    )";

    $pdo->exec($sql);
    echo "Таблица 'players' успешно создана.<br>";
} catch (PDOException $e) {
    die("Ошибка: " . $e->getMessage());
}
echo '<a href="admin_panel.php">Вернуться в админ-панель</a>';
?>