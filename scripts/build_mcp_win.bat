@echo off

echo=
echo=
echo ---------------------------------------------------------------
echo build chromadesk-mcp (Rust MCP Bridge)
echo ---------------------------------------------------------------

:: get script absolute path
set script_path=%~dp0
set old_cd=%cd%
cd /d %~dp0

SETLOCAL EnableDelayedExpansion
set build_mode=release
set errno=1

echo=
echo=
echo ---------------------------------------------------------------
echo parse arguments
echo ---------------------------------------------------------------

:parse_args
if "%1"=="" goto args_done
if /i "%1"=="debug" set build_mode=debug
if /i "%1"=="release" set build_mode=release
shift
goto parse_args
:args_done

echo [*] build mode: %build_mode%
echo=

set mcp_dir=%script_path%..\chromadesk-mcp
set output_path=%script_path%..\output\x64

echo [*] mcp dir: %mcp_dir%
echo [*] output path: %output_path%

:: check if Rust is installed
where cargo >nul 2>nul
if %errorlevel% neq 0 (
    echo [!] error: cargo not found. Please install Rust: https://rustup.rs
    goto return
)

:: build
cd /d "%mcp_dir%"
echo [*] building chromadesk-mcp...

if /i "%build_mode%"=="debug" (
    cargo build
    if not %errorlevel%==0 (
        echo [!] cargo build failed
        goto return
    )
    set cargo_out=%mcp_dir%\target\debug
    set dest_dir=%output_path%\Debug
) else (
    cargo build --release
    if not %errorlevel%==0 (
        echo [!] cargo build failed
        goto return
    )
    set cargo_out=%mcp_dir%\target\release
    set dest_dir=%output_path%\Release
)

:: copy to output directory
if not exist "!dest_dir!" mkdir "!dest_dir!"
copy /Y "!cargo_out!\chromadesk-mcp.exe" "!dest_dir!\" >nul
echo [*] copied chromadesk-mcp.exe to !dest_dir!

echo=
echo=
echo ---------------------------------------------------------------
echo [*] chromadesk-mcp build finished!
echo ---------------------------------------------------------------

set errno=0

:return
cd /d "%old_cd%"
exit /B %errno%

ENDLOCAL
