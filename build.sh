#!/bin/bash

# 编译 Windows
GOOS=windows GOARCH=amd64 go build -o gorsh-windows-amd64.exe main.go

# 编译 Linux
GOOS=linux GOARCH=amd64 go build -o gorsh-linux-amd64 main.go

# 编译 macOS
GOOS=darwin GOARCH=amd64 go build -o gorsh-darwin-amd64 main.go

# 编译 Linux ARM64
GOOS=linux GOARCH=arm64 go build -o gorsh-linux-arm64 main.go