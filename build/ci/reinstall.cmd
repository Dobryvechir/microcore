go install github.com/Dobryvechir/microcore/pkg/main
echo @del %GOPATH%\bin\microcore.exe
echo @ren %GOPATH%\bin\main.exe microcore.exe
echo @copy %GOPATH%\bin\microcore.exe bin\microcore.exe

