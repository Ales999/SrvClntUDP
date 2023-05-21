@echo off

go build -o udpClient.exe Client/main.go
go build -o udpServer.exe Server/main.go

