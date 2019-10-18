@echo off
set GOARCH=amd64
set GOOS=linux
go build -ldflags "-s -w" -o server.out .\server.go
set GOARCH=amd64
set GOOS=windows
go build -ldflags "-s -w" -o server.exe .\server.go