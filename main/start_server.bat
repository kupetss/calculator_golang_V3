@echo off
REM
tasklist /FI "IMAGENAME eq main.exe" 2>NUL | find /I /N "main.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo Сервер уже запущен.
) else (
    echo Запуск сервера...
    start go run main.go
)

REM
timeout /t 4 /nobreak >nul

REM
start http://localhost:8080