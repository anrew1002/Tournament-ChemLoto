<?php
session_start();

if ($_SERVER["REQUEST_METHOD"] == "POST") {
    if (isset($_POST['code'])) {
        $enteredCode = $_POST['code'];
        $correctCode = "123"; // Правильный код организатора

        if ($enteredCode === $correctCode) {
            // Введен правильный код организатора
            if (isset($_SESSION['username'])) {
                // Если организатор уже зарегистрирован, то обновляем его статус
                $dbhost = 'localhost'; // Хост базы данных (в вашем случае, localhost)
            $dbname = 'vloodek20'; // Имя вашей базы данных
            $dbuser = 'root'; // Пользователь базы данных
            $dbpass = ''; // Пароль (пустой, так как у вас нет пароля)

            try {
                $pdo = new PDO("mysql:host=$dbhost;dbname=$dbname", $dbuser, $dbpass);
                    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

                    $name = $_SESSION['username'];

                    // Обновляем статус пользователя на "организатор"
                    $query = "UPDATE players SET is_organizer = 1 WHERE name = :name";
                    $stmt = $pdo->prepare($query);
                    $stmt->bindParam(':name', $name);
                    $stmt->execute();

                    // Обновляем сессию
                    $_SESSION['isOrganizer'] = 1;

                    header('Location: admin_panel.php'); // Перенаправляем на админ-панель
                    exit();
                } catch (PDOException $e) {
                    echo 'Произошла ошибка при обновлении статуса организатора. Пожалуйста, попробуйте позже.';
                }
            } else {
                echo 'Пользователь не вошел в систему. Войдите, чтобы продолжить.';
            }
        } else {
            echo 'Неправильный код организатора.';
        }
    }
}
?>