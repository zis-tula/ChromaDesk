@echo off
echo ---------------------------------------------------------------
echo Build WebClient (Vue 3 + Vite)
echo ---------------------------------------------------------------

set script_path=%~dp0
set webclient_path=%script_path%..\WebClient

cd /d "%webclient_path%"

echo [*] Installing dependencies...
call npm install
if errorlevel 1 (
    echo [!] npm install failed
    exit /B 1
)

echo [*] Building...
call npm run build
if errorlevel 1 (
    echo [!] npm run build failed
    exit /B 1
)

echo [*] Copying remote.html and assets to dist...
xcopy /Y /E /Q "%webclient_path%\js" "%webclient_path%\dist\js\" >nul
copy /Y "%webclient_path%\remote.html" "%webclient_path%\dist\" >nul
copy /Y "%webclient_path%\favicon.ico" "%webclient_path%\dist\" >nul
if exist "%webclient_path%\images" (
    xcopy /Y /E /Q "%webclient_path%\images" "%webclient_path%\dist\images\" >nul
)

echo [*] WebClient build complete: %webclient_path%\dist
exit /B 0
