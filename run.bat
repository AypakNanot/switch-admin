@echo off
cd /d %~dp0
echo Starting Switch Admin Server...
"%~dp0switch-admin.exe" > "%~dp0server.log" 2>&1
echo Server stopped.
