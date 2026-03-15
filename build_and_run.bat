@echo off
cd /d %~dp0
echo Stopping existing processes...
taskkill /F /IM switch-admin.exe 2>nul

echo Building switch-admin...
set CGO_ENABLED=1
go build -o switch-admin.exe ./cmd/main.go
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

echo Starting server...
start "Switch Admin Server" switch-admin.exe
echo Server started! Check http://localhost:9033
