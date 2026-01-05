@echo off
REM KasirNest Build Script for Windows
REM This script builds the KasirNest application with optimization and obfuscation

setlocal enabledelayedexpansion

REM Configuration
set APP_NAME=kasirnest
set VERSION=1.0.0
set BUILD_DIR=build\output
set DIST_DIR=build\dist

REM Build flags
set LDFLAGS=-s -w -X main.version=%VERSION% -X main.buildTime=%DATE%_%TIME%
set CGO_ENABLED=0

echo === KasirNest Build Script ===
echo Version: %VERSION%
echo Build Time: %DATE% %TIME%
echo.

REM Check dependencies
echo Checking dependencies...

REM Check Go version
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo Go version: %GO_VERSION%

REM Check if garble is installed
garble version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Installing garble for obfuscation...
    go install mvdan.cc/garble@latest
)

REM Check if fyne is installed
fyne version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Installing fyne build tools...
    go install fyne.io/fyne/v2/cmd/fyne@latest
)

REM Create build directories
echo Creating build directories...
if not exist %BUILD_DIR% mkdir %BUILD_DIR%
if not exist %DIST_DIR% mkdir %DIST_DIR%

REM Clean previous builds
echo Cleaning previous builds...
del /Q %BUILD_DIR%\* 2>nul
del /Q %DIST_DIR%\* 2>nul

REM Download dependencies
echo Downloading dependencies...
go mod download
go mod tidy

REM Parse command line arguments
set OBFUSCATE=true
set BUILD_ALL=false
set RUN_TESTS=true

:parse_args
if "%~1"=="" goto args_done
if "%~1"=="--no-obfuscate" (
    set OBFUSCATE=false
    shift
    goto parse_args
)
if "%~1"=="--all-platforms" (
    set BUILD_ALL=true
    shift
    goto parse_args
)
if "%~1"=="--no-test" (
    set RUN_TESTS=false
    shift
    goto parse_args
)
echo Unknown option: %~1
echo Usage: %0 [--no-obfuscate] [--all-platforms] [--no-test]
exit /b 1

:args_done

REM Run tests (optional)
if "%RUN_TESTS%"=="true" (
    echo Running tests...
    go test ./... -v
    if %ERRORLEVEL% NEQ 0 (
        echo Tests failed!
        exit /b 1
    )
)

REM Build function - Windows
:build_windows
set PLATFORM=%1
set ARCH=%2
set OUTPUT=%BUILD_DIR%\%APP_NAME%_%PLATFORM%_%ARCH%.exe

echo Building for %PLATFORM%/%ARCH%...

set GOOS=%PLATFORM%
set GOARCH=%ARCH%
set CGO_ENABLED=%CGO_ENABLED%

REM Build with or without obfuscation
if "%OBFUSCATE%"=="true" (
    echo Building with obfuscation...
    garble -literals -seed=random build -ldflags="%LDFLAGS%" -o "%OUTPUT%" .
) else (
    echo Building without obfuscation...
    go build -ldflags="%LDFLAGS%" -o "%OUTPUT%" .
)

if %ERRORLEVEL% EQU 0 (
    echo ✓ Built: %OUTPUT%
    for %%A in ("%OUTPUT%") do echo   Size: %%~zA bytes
) else (
    echo ✗ Failed to build for %PLATFORM%/%ARCH%
    exit /b 1
)
goto :eof

REM Build function - Unix
:build_unix
set PLATFORM=%1
set ARCH=%2
set OUTPUT=%BUILD_DIR%\%APP_NAME%_%PLATFORM%_%ARCH%

echo Building for %PLATFORM%/%ARCH%...

set GOOS=%PLATFORM%
set GOARCH=%ARCH%
set CGO_ENABLED=%CGO_ENABLED%

REM Build with or without obfuscation
if "%OBFUSCATE%"=="true" (
    echo Building with obfuscation...
    garble -literals -seed=random build -ldflags="%LDFLAGS%" -o "%OUTPUT%" .
) else (
    echo Building without obfuscation...
    go build -ldflags="%LDFLAGS%" -o "%OUTPUT%" .
)

if %ERRORLEVEL% EQU 0 (
    echo ✓ Built: %OUTPUT%
    for %%A in ("%OUTPUT%") do echo   Size: %%~zA bytes
) else (
    echo ✗ Failed to build for %PLATFORM%/%ARCH%
    exit /b 1
)
goto :eof

REM Start builds
echo Starting builds...

if "%BUILD_ALL%"=="true" (
    REM Build for all platforms
    call :build_windows windows amd64
    call :build_windows windows 386
    call :build_unix linux amd64
    call :build_unix linux 386
    call :build_unix darwin amd64
    call :build_unix darwin arm64
) else (
    REM Build for Windows only
    call :build_windows windows amd64
)

REM Create distribution packages
echo Creating distribution packages...

for %%F in (%BUILD_DIR%\*) do (
    set filename=%%~nxF
    set platform_arch=!filename:%APP_NAME%_=!
    
    echo Packaging !filename!...
    
    REM Create package directory
    set package_dir=%DIST_DIR%\%APP_NAME%_!platform_arch!
    if not exist "!package_dir!" mkdir "!package_dir!"
    
    REM Copy binary
    copy "%%F" "!package_dir!\"
    
    REM Copy configuration template
    copy "config\app.ini.example" "!package_dir!\"
    
    REM Copy assets
    xcopy "assets" "!package_dir!\assets\" /E /I /Q >nul 2>&1
    
    REM Copy documentation
    copy "README.md" "!package_dir!\" >nul 2>&1
    copy "FIREBASE_SETUP.md" "!package_dir!\" >nul 2>&1
    
    REM Create ZIP archive
    cd %DIST_DIR%
    powershell -command "Compress-Archive -Path '%APP_NAME%_!platform_arch!' -DestinationPath '%APP_NAME%_!platform_arch!.zip' -Force"
    cd ..
    
    echo ✓ Package created: %APP_NAME%_!platform_arch!.zip
)

REM Build summary
echo.
echo === Build Summary ===
echo Built binaries:
dir %BUILD_DIR%

echo.
echo Distribution packages:
dir %DIST_DIR%\*.zip

echo.
echo ✓ Build completed successfully!

REM Optional: Run the binary
if "%~1"=="--run" (
    echo Starting application...
    %BUILD_DIR%\%APP_NAME%_windows_amd64.exe
)

endlocal