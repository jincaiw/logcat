@echo off
echo Building logcat Web Server...

cd /d "%~dp0"

echo 1. Installing frontend dependencies...
cd frontend
call npm install
cd ..

echo 2. Type checking frontend...
cd frontend
call npm run typecheck
cd ..

echo 3. Building frontend...
cd frontend
call npm run build
cd ..

echo 4. Building web server binary...
go build -o build\bin\logcat-web.exe

echo 5. Copying templates...
if exist templates (
    xcopy /E /I /Y templates build\bin\templates
)

echo Done! Output: build/bin/
echo.
echo Files:
echo   - logcat-web.exe (binary)
echo   - templates/ (parse_templates.json, filter_policies.json)
echo.
echo Usage: logcat-web.exe [port]
echo Default port: 8080

pause
