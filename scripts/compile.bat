@Echo Off

set EXECUTABLE=binary\giggle.exe
set CONTENT_FOLDER=gigglecontent
set OUT_FILE=content.go

if exist del /F %EXECUTABLE%
if exist del /F %CONTENT_FOLDER%\%OUT_FILE%

REM go format the code
gofmt -w .

REM convert icons into go code
cd %CONTENT_FOLDER%
go-bindata.exe -o %OUT_FILE% -pkg %CONTENT_FOLDER% index.html css/... images/... js/...
cd ..

REM build & install
go install -v
go build -o %EXECUTABLE%
