@Echo Off

set EXECUTABLE=binary\giggle.exe
set SERVICE_FOLDER=giggleservice
set ICONS_FILE=icons.go

if exist del /F %EXECUTABLE%
if exist del /F %SERVICE_FOLDER%\%ICONS_FILE%

REM go format the code
gofmt -w .

REM convert icons into go code
cd %SERVICE_FOLDER%
go-bindata.exe -o %ICONS_FILE% -pkg %SERVICE_FOLDER% icons
cd ..\

REM build & install
go install
go build -o %EXECUTABLE%
