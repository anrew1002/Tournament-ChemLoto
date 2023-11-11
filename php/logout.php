<?php
// Начнем сессию, если она еще не начата
if (session_status() == PHP_SESSION_NONE) {
    session_start();
}

// Уничтожаем все данные сессии
session_unset();

// Завершаем текущую сессию
session_destroy();

// Перенаправляем пользователя на страницу входа (index.php)
header('Location: index.php');
?>