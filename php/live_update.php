<?php
header('Content-Type: text/event-stream');
header('Cache-Control: no-cache');

while (true) {
    $randomNumber = rand(1, 100);
    echo "data: " . json_encode(['randomNumber' => $randomNumber]) . "\n\n";
    ob_flush();
    flush();
    sleep(1);
}
?>