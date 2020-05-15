@Echo Off

set CONTENT_FOLDER=content
set OUT_FILE=content.go

REM cleaning up
if exist del /F %CONTENT_FOLDER%\%OUT_FILE%

REM install go-bindata
go get -u github.com/go-bindata/go-bindata/...

REM convert icons into go code
cd %CONTENT_FOLDER%
go-bindata.exe -o %OUT_FILE% -pkg %CONTENT_FOLDER% images/...
cd ..

REM build & install
go mod tidy
go install -v
