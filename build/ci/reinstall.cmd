go install github.com/Dobryvechir/microcore/pkg/main
@del %GOPATH%\bin\microcore.exe
@ren %GOPATH%\bin\main.exe microcore.exe
@copy %GOPATH%\bin\microcore.exe C:\prg\go\bin\microcore.exe

