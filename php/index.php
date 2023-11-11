<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f5f5f5;
        }

        #header {
            text-align: right;
            padding: 10px;
        }

        #content {
            max-width: 400px;
            margin: 0 auto;
            background-color: #fff;
            border: 1px solid #ccc;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }

        #content h2 {
            text-align: center;
        }

        #content label {
            display: block;
            margin-bottom: 10px;
        }

        #content input[type="text"],
        #content input[type="password"] {
            width: 100%;
            padding: 10px;
            margin-bottom: 15px; /* Увеличил отступ между инпутами */
            border: 1px solid #ccc;
            border-radius: 5px;
        }

        #content input[type="checkbox"] {
            margin-right: 5px;
        }

        #content button {
            display: block;
            width: 100%;
            padding: 10px;
            margin-bottom: 15px; /* Увеличил отступ между кнопками */
            background-color: #007BFF;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }

        .hidden {
            display: none;
        }

        #toggleForm {
            display: block;
            width: 100%;
            padding: 10px;
            background-color: #28a745;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }

        #toggleForm:hover {
            background-color: #218838;
        }

        #error {
            color: red;
            text-align: center;
            margin-top: 10px;
        }
    </style>
</head>
<body>
<?php
session_start();

// Проверка, если пользователь уже вошел в систему, перенаправляем его на страницу с комнатами игры
if (isset($_SESSION['username'])) {
    header('Location: room_list.php');
    exit();
}

$errorMessage = ''; // Переменная для хранения сообщений об ошибках

if ($_SERVER["REQUEST_METHOD"] == "POST") {
    if (isset($_POST['register'])) {
        // Обработка регистрации
        $name = $_POST["name"];
        $password = $_POST["password"];
        $isOrganizer = isset($_POST["organizer"]);
        $code = "123"; // Код для проверки
        if ($isOrganizer) {
            $enteredCode = $_POST['code'];

            if ($enteredCode !== $code) {
                $errorMessage = 'Неправильный код организатора.';

            }
        }
        // Простая проверка, что поля не пустые
        if (empty($name) || empty($password) || ($isOrganizer && empty($code))) {
            $errorMessage = 'Пожалуйста, заполните все обязательные поля.';
        }
        elseif  ($isOrganizer && ($enteredCode !==$code) ){
            $errorMessage = 'Неправильный код организатора.';
        }
         else {
            // Подключение к базе данных и проверка наличия пользователя с таким именем
            // $dbname = "vloodek20";
            // $firstname = "vloodek20";
            // $dbpassword = "2IRj7JgP";
            // $servername = "127.0.0.1";
            $dbhost = 'localhost'; // Хост базы данных (в вашем случае, localhost)
            $dbname = 'vloodek20'; // Имя вашей базы данных
            $dbuser = 'root'; // Пользователь базы данных
            $dbpass = ''; // Пароль (пустой, так как у вас нет пароля)

            try {
                $pdo = new PDO("mysql:host=$dbhost;dbname=$dbname", $dbuser, $dbpass);
                $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

                // Проверка существования пользователя с таким же именем
                $query = "SELECT * FROM players WHERE name = :name";
                $stmt = $pdo->prepare($query);
                $stmt->bindParam(':name', $name);
                $stmt->execute();

                if ($stmt->rowCount() > 0) {
                    $errorMessage = 'Пользователь с таким именем уже существует.';
                } else {
                    // Если пользователя с таким именем нет, создаем новую запись
                    $query = "INSERT INTO players (name, password, is_organizer) VALUES (:name, :password, :is_organizer)";
                    $stmt = $pdo->prepare($query);
                    $stmt->bindParam(':name', $name);
                    $stmt->bindParam(':password', $password);
                    $stmt->bindParam(':is_organizer', $isOrganizer, PDO::PARAM_INT);
                    $stmt->execute();

                    // Устанавливаем сессию
                    $_SESSION['username'] = $name;
                    $_SESSION['isOrganizer'] = $isOrganizer;

                    header('Location: room_list.php');
                    exit();
                }
            } catch (PDOException $e) {
                $errorMessage = 'Произошла ошибка при регистрации. Пожалуйста, попробуйте позже. ' . $e;
            }
        }
    } elseif (isset($_POST['login'])) {
        // Обработка входа (пользователь уже зарегистрирован)
        $name = $_POST["loginName"];
        $loginPassword = $_POST["loginPassword"];

        // Подключение к базе данных и проверка пароля
        $dbhost = 'localhost'; // Хост базы данных (в вашем случае, localhost)
            $dbname = 'vloodek20'; // Имя вашей базы данных
            $dbuser = 'root'; // Пользователь базы данных
            $dbpass = ''; // Пароль (пустой, так как у вас нет пароля)

            try {
                $pdo = new PDO("mysql:host=$dbhost;dbname=$dbname", $dbuser, $dbpass);
            $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

            // Проверка пользователя по имени и паролю
            $query = "SELECT * FROM players WHERE name = :name AND password = :password";
            $stmt = $pdo->prepare($query);
            $stmt->bindParam(':name', $name);
            $stmt->bindParam(':password', $loginPassword);
            $stmt->execute();

            if ($stmt->rowCount() > 0) {
                // Пользователь с таким именем и паролем существует
                // Устанавливаем сессию
                $_SESSION['username'] = $name;
                $_SESSION['isOrganizer'] = $stmt->fetch(PDO::FETCH_ASSOC)['is_organizer'];
                header('Location: room_list.php');
                exit();
            } else {
                $errorMessage = 'Неправильное имя пользователя или пароль.';
            }
        } catch (PDOException $e) {
            $errorMessage = 'Произошла ошибка при входе. Пожалуйста, попробуйте позже. ' . $e;
        }
    }
}
?>
    
<div id="header">
    <?php
    // Показать имя пользователя, если вошли в систему
    if (isset($_SESSION['username'])) {
        echo 'Здравствуйте, ' . $_SESSION['username'] . '!';
    }
    ?>
</div>

<div id="content">
    <h2>Регистрация</h2>
    <form method="post" action="<?php echo $_SERVER['PHP_SELF']; ?>">
        <label for="name">Введите имя:</label>
        <input type="text" name="name" id="name" required>
        
        <label for="password">Пароль:</label>
        <input type="password" name="password" id="password" required>
        
        <input type="checkbox" name="organizer" id="organizer">
        <label for="organizer">Вы организатор?</label>
        
        <div class="hidden" id="codeContainer">
            <label for="code">Введите код:</label>
            <input type="text" name="code" id="code">
        </div>
        
        <button type="submit" name="register">Зарегистрировать</button>
    </form>

    <!-- Форма для входа -->
    <form method="post" action="<?php echo $_SERVER['PHP_SELF']; ?>">
        <h2>Вход</h2>
        <label for="loginName">Имя пользователя:</label>
        <input type="text" name="loginName" id="loginName" required>
        <label for="loginPassword">Пароль:</label>
        <input type="password" name="loginPassword" id="loginPassword" required>
        <button type="submit" name="login">Войти</button>
    </form>
    
    <div id="error">
        <?php
        // Вывод сообщений об ошибках, если они есть
        if (!empty($errorMessage)) {
            echo $errorMessage;
        }
        ?>
    </div>
</div>

<script>
    // JavaScript для показа/скрытия поля ввода кода при выборе чекбокса
    const organizerCheckbox = document.getElementById('organizer');
    const codeContainer = document.getElementById('codeContainer');

    organizerCheckbox.addEventListener('change', function() {
        codeContainer.style.display = organizerCheckbox.checked ? 'block' : 'none';
    });


</script>

</body>
</html>