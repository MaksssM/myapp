@echo off

REM Создание лог-файла
set LOGFILE=start-all.log
if exist %LOGFILE% del %LOGFILE%
echo Лог запуска сервисов > %LOGFILE%

REM Запуск бэкенда
cd backend
call go run main.go 2>> ..\%LOGFILE%
cd ..

REM Запуск мобильного приложения
cd TipaTwitterMobile
call npx react-native run-android 2>> ..\%LOGFILE%
cd ..

REM Запуск фронтенда
cd frontend
call npm start 2>> ..\%LOGFILE%
cd ..

REM Запуск Centrifugo
cd centrifugo
call centrifugo --config=config.json --admin 2>> ..\%LOGFILE%
cd ..

REM Ожидание завершения
pause
