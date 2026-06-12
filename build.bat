@echo off
REM Quick build and install script for Windows

echo Building devgitsecops...
go build -v -o bin\devgitsecops.exe main.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Build successful! Binary created at: bin\devgitsecops.exe
    echo.
    echo Quick start:
    echo   1. Check tool status: .\bin\devgitsecops.exe status
    echo   2. Install tools:     .\bin\devgitsecops.exe install --all --auto
    echo   3. Use tools:         .\bin\devgitsecops.exe kubectl get pods
    echo.
    echo To add to PATH, run as Administrator:
    echo   setx /M PATH "%%PATH%%;%CD%\bin"
) else (
    echo Build failed!
    exit /b 1
)
