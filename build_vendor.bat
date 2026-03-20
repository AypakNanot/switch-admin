@echo off
chcp 65001 >nul
cd /d %~dp0

echo ========================================
echo Build with vendor mode (offline)
echo ========================================
echo.

echo Setting environment variables...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

echo Cleaning old build files...
del /f /q switch-admin.exe 2>nul

echo Starting compilation (vendor mode)...
go build -mod=vendor -ldflags="-s -w" -o switch-admin.exe ./cmd/main.go

if errorlevel 1 (
    echo.
    echo [FAILED] Build failed!
    echo Tip: Ensure vendor directory exists and is complete
    exit /b 1
)

echo.
echo [SUCCESS] Build completed!
dir switch-admin.exe
echo.
echo Vendor mode works without network connection.
