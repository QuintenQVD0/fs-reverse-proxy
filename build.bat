@echo off
:: Set the output file name
set OUTPUT_FILE=reverse_proxy_linux_x64

:: Check if Go is installed
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH.
    exit /b 1
)

:: Set the environment variables for Linux x64 build
set GOOS=linux
set GOARCH=amd64

:: Compile the Go program
echo Compiling reverse proxy for Linux x64...
go build .\cmd\main.go

rename main %OUTPUT_FILE%

:: Check if the build succeeded
if %errorlevel% neq 0 (
    echo Build failed.
    exit /b 1
) else (
    echo Build succeeded. Output file: %OUTPUT_FILE%
)