@echo off
chcp 65001 >nul
cd /d %~dp0

echo ========================================
echo 使用 vendor 模式编译 (离线开发)
echo ========================================
echo.

echo 设置环境变量...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

echo 清理旧的构建文件...
del /f /q switch-admin.exe 2>nul

echo 开始编译 (vendor 模式)...
go build -mod=vendor -ldflags="-s -w" -o switch-admin.exe ./cmd/main.go

if errorlevel 1 (
    echo.
    echo [失败] 编译失败!
    echo 提示：确保 vendor 目录存在且完整
    exit /b 1
)

echo.
echo [成功] 编译完成!
dir switch-admin.exe
echo.
echo vendor 模式可以在没有网络连接的情况下编译
