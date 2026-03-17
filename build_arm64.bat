@echo off
chcp 65001 >nul
cd /d %~dp0

echo ========================================
echo 交叉编译 switch-admin (Linux ARM64)
echo ========================================
echo.

rem 动态查找 Go 路径
where go >nul 2>&1
if errorlevel 1 (
    echo [错误] 未找到 Go，请确保已安装并添加到 PATH
    exit /b 1
)

echo 设置环境变量...
set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=0

echo 清理旧的构建文件...
del /f /q switch-admin-arm64 2>nul
del /f /q switch-admin-arm64.exe 2>nul

echo 开始编译...
echo   GOOS=%GOOS%
echo   GOARCH=%GOARCH%
echo   CGO_ENABLED=%CGO_ENABLED%
echo.

go build -mod=vendor -ldflags="-s -w" -o switch-admin-arm64 ./cmd/main.go

if errorlevel 1 (
    echo.
    echo [失败] 编译失败!
    exit /b 1
)

echo.
echo [成功] 编译完成!
echo 输出文件：switch-admin-arm64
dir switch-admin-arm64
echo.
echo 提示：将文件上传到交换机的 /root 目录运行
echo 例如：scp switch-admin-arm64 root@switch:/root/
echo.
