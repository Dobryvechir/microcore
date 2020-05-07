go install github.com/Dobryvechir/microcore/src/mainParse
@del %GOPATH%\bin\microcorep.exe
@ren %GOPATH%\bin\mainParse.exe microcorep.exe
@copy %GOPATH%\bin\microcorep.exe src\mainParse\microcorep.exe
cd src\mainParse\
microcorep.exe

