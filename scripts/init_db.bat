@echo off
REM Concurrent Knowledge - Database Initialization Script
REM This script creates the SQLite database and runs migrations

echo ========================================
echo Concurrent Knowledge - DB Initialization
echo ========================================
echo.

REM Set database path
set DB_PATH=..\data\quiz.db

REM Check if data directory exists, if not create it
if not exist ..\data mkdir ..\data

REM Check if database already exists
if exist %DB_PATH% (
    echo WARNING: Database already exists at %DB_PATH%
    set /p CONFIRM="Do you want to delete and recreate it? (y/n): "
    if /i "%CONFIRM%" neq "y" (
        echo Aborted.
        exit /b 1
    )
    echo Deleting existing database...
    del %DB_PATH%
)

echo Creating new database...
echo.

REM Run schema migration
echo Executing schema migration (001_schema.sql)...
sqlite3 %DB_PATH% < ..\migrations\001_schema.sql
if %ERRORLEVEL% neq 0 (
    echo ERROR: Failed to execute schema migration
    exit /b 1
)
echo Schema migration completed successfully.
echo.

REM Run seed data migration
echo Executing seed data (002_seed.sql)...
sqlite3 %DB_PATH% < ..\migrations\002_seed.sql
if %ERRORLEVEL% neq 0 (
    echo ERROR: Failed to execute seed data
    exit /b 1
)
echo Seed data loaded successfully.
echo.

REM Verify database
echo Verifying database...
sqlite3 %DB_PATH% "SELECT COUNT(*) as genre_count FROM genres;"
sqlite3 %DB_PATH% "SELECT COUNT(*) as question_count FROM questions;"
echo.

echo ========================================
echo Database initialization completed!
echo ========================================
echo Database location: %DB_PATH%
echo.

pause
