@echo off
chcp 65001 >nul
echo ========================================
echo    switch-admin 编译脚本（无 CGO）
echo ========================================
echo.
echo 使用纯 Go SQLite 驱动 (modernc.org/sqlite)
echo 不需要安装 GCC，编译后的程序可在任意 Windows 电脑运行
echo.

echo [1/3] 设置无 CGO 环境...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

echo [2/3] 开始编译...
echo 输出文件：switch-admin.exe
go build -mod=vendor -o switch-admin.exe .\cmd

if %ERRORLEVEL% EQU 0 (
    echo.
    echo [3/3] 编译成功!
    echo.
    echo ========================================
    echo   构建完成（无 CGO 版本）
    echo ========================================
    echo.
    echo 文件信息:
    for %%A in (switch-admin.exe) do (
        echo   文件名：%%A
        echo   大小：%%~zA 字节
    )
    echo.
    echo 下一步操作:
    echo 1. 将 switch-admin.exe 复制到目标电脑
    echo 2. 确保目标电脑有 data/ 和 uploads/ 目录
    echo 3. 运行：switch-admin.exe
    echo.
    echo 优点：
    echo - 不需要 GCC 环境
    echo - 编译后的程序可在任意 Windows 电脑运行
    echo - 文件更小（约 50MB）
    echo.
) else (
    echo.
    echo ========================================
    echo   编译失败!
    echo ========================================
    echo.
    echo 请检查错误信息
    echo.
)

pause
