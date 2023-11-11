<?php
session_start();

if (!isset($_SESSION['username'])) {
    header('Location: index.php'); // Перенаправление на index.php, если пользователь не вошел
    exit();
}

if (isset($_POST['clearRooms'])) {
    // Подключение к базе данных (замените значения на свои)
    $dbhost = 'localhost'; // Хост базы данных (в вашем случае, localhost)
            $dbname = 'vloodek20'; // Имя вашей базы данных
            $dbuser = 'root'; // Пользователь базы данных
            $dbpass = ''; // Пароль (пустой, так как у вас нет пароля)

            try {
                $pdo = new PDO("mysql:host=$dbhost;dbname=$dbname", $dbuser, $dbpass);
        $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

        // Удалить все комнаты из базы данных
        $query = "DELETE FROM rooms";
        $pdo->exec($query);

        // Очистить папку "rooms" от файлов
        $roomsDirectory = 'rooms';
        if (is_dir($roomsDirectory)) {
            $files = glob($roomsDirectory . '/*');
            foreach ($files as $file) {
                if (is_file($file)) {
                    unlink($file);
                }
            }
        }
    } catch (PDOException $e) {
        echo 'Ошибка при очистке комнат: ' . $e->getMessage();
    }
}

if (isset($_POST['clearRooms'])) {
    // Подключение к базе данных (замените значения на свои)
    $dbname = "vloodek20";
    $firstname = "vloodek20";
    $password = "2IRj7JgP";
    $servername = "127.0.0.1";

    try {
        $pdo = new PDO("mysql:host=$servername;dbname=$dbname", $firstname, $password);
        $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

        // Удалить таблицу "rooms" (если она существует)
        $query = "DROP TABLE IF EXISTS rooms";
        $pdo->exec($query);

        // Создать таблицу "rooms" заново
        $query = "
            CREATE TABLE rooms (
                id INT AUTO_INCREMENT PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                is_auto TINYINT(1) DEFAULT 0,
                time INT DEFAULT 0,
                max_players INT
            );
        ";
        $pdo->exec($query);
    } catch (PDOException $e) {
        echo 'Ошибка при пересоздании таблицы rooms: ' . $e->getMessage();
    }
}
?>

<!DOCTYPE html>
<html>
<head>
    <title>Админ-панель</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f5f5f5;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
        }

        .container {
            max-width: 800px;
            background-color: #fff;
            border: 1px solid #ccc;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            text-align: center;
        }

        h1 {
            font-size: 24px;
            margin-bottom: 20px;
        }

        p {
            font-size: 18px;
            margin-bottom: 20px;
        }

        .button-container {
            display: flex;
            justify-content: center;
            align-items: center;
            flex-direction: column;
        }

        .button {
            display: inline-block;
            padding: 10px 20px;
            background-color: #007BFF;
            color: #fff;
            text-decoration: none;
            font-size: 18px;
            margin: 5px;
            border-radius: 5px;
            cursor: pointer;
            transition: background-color 0.3s;
        }

        .button.red {
            background-color: #ff0000;
        }

        input[type="text"] {
            font-size: 18px;
            padding: 10px;
            border: 1px solid #ccc;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <h1>Админ-панель</h1>
    <div class="container">
        <?php
        if (isset($_SESSION['isOrganizer']) && $_SESSION['isOrganizer']) {
            // Организатор вошел в систему, показываем кнопку для очистки комнат и пользователей
            echo '<p>Здравствуйте, ' . $_SESSION['username'] . '!</p>';
            echo '<div class="button-container">';
            echo '<form method="post" action="admin_panel.php">';
            echo '<button type="submit" name="clearRooms" class="button red">Очистить список комнат</button>';
            echo '</form>';
            echo '<form method="post" action="admin_panel.php">';
            echo '<button type="submit" name="clearUsers" class="button red">Очистить базу данных пользователей</button>';
            echo '</form>';
            echo '</div>';
            echo '<br>';
            echo '<a href="room_list.php" class="button">Вернуться к комнатам</a>';
        } else {
            // Организатор не вошел в систему, показываем форму для ввода кода
            echo '
            <form method="post" action="login_organizer.php">
                <label for="code">Введите код:</label>
                <input type="text" name="code" id="code">
                <button type="submit" class="button">Войти как организатор</button>
            </form>';
        }
        ?>
    </div>
</body>
</html>
