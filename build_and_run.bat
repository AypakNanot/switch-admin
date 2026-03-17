@echo off
cd /d %~dp0

rem 动态查找 Go 路径
where go >nul 2>&1
if errorlevel 1 (
    echo [错误] 未找到 Go，请确保已安装并添加到 PATH
    exit /b 1
)

echo Stopping existing processes...
taskkill /F /IM switch-admin.exe 2>nul

echo Cleaning old build files...
del /f /q switch-admin.exe 2>nul
del /f /q gin-admin.exe 2>nul
del /f /q switch-admin-*.exe 2>nul
if exist dist rmdir /s /q dist 2>nul
if exist build rmdir /s /q build 2>nul
if exist data\bak rmdir /s /q data\bak 2>nul

echo Building switch-admin...
set CGO_ENABLED=0
set GOARCH=amd64
go build -mod=vendor -ldflags="-s -w" -o switch-admin.exe ./cmd/main.go
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

echo Starting server...
start "Switch Admin Server" switch-admin.exe
echo Server started! Check http://localhost:9033
