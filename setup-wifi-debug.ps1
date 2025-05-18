# Скрипт для настройки отладки по Wi-Fi

# Установка модуля для генерации QR-кодов
if (-not (Get-Module -ListAvailable -Name QRCodeGenerator)) {
    Install-Module -Name QRCodeGenerator -Force -Scope CurrentUser
}

# Переключение устройства в режим Wi-Fi-отладки
Write-Host "Переключение устройства в режим Wi-Fi-отладки..."
adb tcpip 5555

# Запрос IP-адреса у пользователя
$ipAddress = Read-Host "Введите IP-адрес устройства (например, 192.168.1.100)"

# Подключение к устройству по Wi-Fi
Write-Host "Подключение к устройству по IP-адресу $ipAddress..."
adb connect "$ipAddress:5555"

# Генерация QR-кода для подключения
Write-Host "Генерация QR-кода для подключения..."
$qrContent = "adb connect $ipAddress:5555"
$qrPath = "wifi-debug-qr.png"
New-QRCode -Content $qrContent -Path $qrPath

Write-Host "QR-код сохранён в $qrPath. Отсканируйте его для подключения."

# Проверка подключения
Write-Host "Проверка подключенных устройств..."
adb devices

Write-Host "Настройка завершена. Если устройство не отображается, проверьте настройки и повторите попытку."
