@echo off

go build -ldflags="-s" -o udpClient.exe Client/udpClient.go
go build -ldflags="-s" -o udpServer.exe Server/udpServer.go


