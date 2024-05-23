@echo off
SET VERSION=1.1.0

REM Create the version directory under release
mkdir release\%VERSION%

REM Build for Windows 64 bit
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-s -w" -o release\%VERSION%\shared2-%VERSION%-windows-amd64.exe

REM Build for Windows 32 bit
SET GOOS=windows
SET GOARCH=386
go build -ldflags="-s -w" -o release\%VERSION%\shared2-%VERSION%-windows-386.exe

REM Build for Mac
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags="-s -w" -o release\%VERSION%\shared2-%VERSION%-darwin-amd64

REM Build for Linux amd64
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags="-s -w" -o release\%VERSION%\shared2-%VERSION%-linux-amd64

REM Build for Linux arm
SET GOOS=linux
SET GOARCH=arm
go build -ldflags="-s -w" -o release\%VERSION%\shared2-%VERSION%-linux-arm
