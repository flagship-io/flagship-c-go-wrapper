#!/bin/bash

# downloading dependencies
echo "downloading dependencies"
go mod download

echo "building linux library"
GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o build/libflagship_linux_amd64.so

echo "building mac library"
GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ go build -buildmode=c-shared -o build/libflagship_darwin_amd64.dylib

echo "building windows library"
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -buildmode=c-shared -o build/libflagship_windows_amd64.dll