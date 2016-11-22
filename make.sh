#!/usr/bin/env bash

version=$(grep -F "VERSION = " settings.go | cut -d\" -f2)

echo "Build frontend app"

cd admin_frontend
npm run build
go-bindata -pkg core -prefix dist -o ../admin_frontend.go dist/admin/...
cd -

echo "Cross compiling gopypi version: $version"

echo "Compiling for linux-amd64..."
env GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o build/gopypi-linux-amd64-$version ./gopypi
echo "Compiling for linux-arm64..."
env GOOS=linux GOARCH=arm64 go build -ldflags "-s" -o build/gopypi-linux-arm64-$version ./gopypi
echo "Compiling for darwin-amd64..."
env GOOS=darwin GOARCH=amd64 go build -ldflags "-s" -o build/gopypi-darwin-amd64-$version ./gopypi
echo "Compiling for freebsd-amd64..."
env GOOS=freebsd GOARCH=amd64 go build -ldflags "-s" -o build/gopypi-freebsd-amd64-$version ./gopypi
echo "Compiling for windows-amd64..."
env GOOS=windows GOARCH=amd64 go build -ldflags "-s" -o build/gopypi-windows-amd64-$version.exe ./gopypi
