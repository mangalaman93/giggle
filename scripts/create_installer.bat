@Echo Off

set INSTALLER_FOLDER=binary
set W_INSTALLER=giggle.msi
set TEMPLATES_FOLDER=templates
set INSTALLER_FILE=%INSTALLER_FOLDER%\%W_INSTALLER%

if [%1]==[] (
  echo version number required
  exit /b 1
)

REM clean up
if exist %INSTALLER_FILE% del /F %INSTALLER_FILE%

REM compile
call scripts/compile.bat

go-msi check-json
go-msi set-guid
if not exist %INSTALLER_FOLDER% mkdir %INSTALLER_FOLDER%
go-msi make --src %TEMPLATES_FOLDER% --version %1 --msi %INSTALLER_File%
