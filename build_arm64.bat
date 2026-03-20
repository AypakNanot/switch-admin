@echo off
chcp 65001 >nul
cd /d %~dp0

echo ========================================
echo Cross-compile switch-admin (Linux ARM64)
echo ========================================
echo.

rem Find Go path dynamically
where go >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go not found. Please ensure Go is installed and added to PATH.
    exit /b 1
)

echo Setting environment variables...
set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=0

echo Cleaning old build files...
del /f /q switch-admin-arm64 2>nul
del /f /q switch-admin-arm64.exe 2>nul

echo Starting compilation...
echo   GOOS=%GOOS%
echo   GOARCH=%GOARCH%
echo   CGO_ENABLED=%CGO_ENABLED%
echo.

go build -mod=vendor -ldflags="-s -w" -o switch-admin-arm64 ./cmd/main.go

if errorlevel 1 (
    echo.
    echo [FAILED] Build failed!
    exit /b 1
)

echo.
echo [SUCCESS] Build completed!
echo Output file: switch-admin-arm64
dir switch-admin-arm64
echo.
echo Tip: Upload to switch /root directory to run
echo Example: scp switch-admin-arm64 root@switch:/root/
echo.
