@echo off

set BINARY_NAME=luidium.exe

set CURRENT_DIR=%~dp0

set BINARY_PATH=%CURRENT_DIR%luidium-windows-amd64.exe

set LOCAL_BIN_DIR=%LOCALAPPDATA%\bin

if not exist "%LOCAL_BIN_DIR%" mkdir "%LOCAL_BIN_DIR%"

copy "%BINARY_PATH%" "%LOCAL_BIN_DIR%\%BINARY_NAME%"

setx PATH "%PATH%;%LOCAL_BIN_DIR%"

echo Luidium CLI installed successfully