<?php
require 'vendor/autoload.php';

use Ratchet\ConnectionInterface;
use Ratchet\Http\HttpServer;
use Ratchet\Server\IoServer;
use Ratchet\WebSocket\WsServer;
use Ratchet\MessageComponentInterface;
use React\EventLoop\Factory;

class MyWebSocketServer implements MessageComponentInterface {
    protected $clients;
    protected $loop;

    public function __construct() {
        $this->clients = new \SplObjectStorage;
        $this->loop = Factory::create();
    }

    public function onOpen(ConnectionInterface $conn) {
        $this->clients->attach($conn);
        echo "Client connected ({$conn->resourceId})\n";
    }

    public function onMessage(ConnectionInterface $from, $msg) {
        // Ничего не делаем при получении сообщения от клиента
    }

    public function onClose(ConnectionInterface $conn) {
        $this->clients->detach($conn);
        echo "Client disconnected ({$conn->resourceId})\n";
    }

    public function onError(ConnectionInterface $conn, \Exception $e) {
        echo "An error occurred: {$e->getMessage()}\n";
        $conn->close();
    }

    public function sendRandomNumber() {
    $randomNumber = rand(1, 100); // Генерируем случайное число от 1 до 100
    foreach ($this->clients as $client) {
        $client->send(json_encode(['randomNumber' => $randomNumber]));
    }
    echo "Generated random number: {$randomNumber}\n"; // Добавьте эту строку для отладки
}

public function run() {
    $server = IoServer::factory(new HttpServer(new WsServer($this)), 8080);

    $this->loop = $server->loop; // Используйте встроенный Event Loop

    echo "Server started\n";

    $this->loop->addPeriodicTimer(5, function () {
        $this->sendRandomNumber();
    });

    $server->run();
}
}

$server = new MyWebSocketServer();
$server->run();
