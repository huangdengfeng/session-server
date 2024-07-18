#!/bin/bash
set -e
go test ./...
#go test -bench=. ./...
DIR=$(cd "$(dirname $0)" && pwd)
APP_NAME="${DIR##*/}"
cd $DIR
rm -rf target
mkdir -p target/$APP_NAME
\cp -R bin target/$APP_NAME
\cp -R conf target/$APP_NAME
# -o 后加目录则放入目录，不是目录则为产出物名称
# 使用uname命令来判断操作系统类型
OS=$(uname -s)

case "$OS" in
    Linux*)
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o target/$APP_NAME/bin/$APP_NAME cmd/*
        ;;
    Darwin*)
        CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o target/$APP_NAME/bin/$APP_NAME cmd/*
        ;;
    CYGWIN*|MINGW32*|MSYS*|MINGW*)
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o target/$APP_NAME/bin/$APP_NAME cmd/*
        ;;
    *)
        echo "Unknown operating system"
        exit 1
        ;;
esac
go build -o target/$APP_NAME/bin/$APP_NAME cmd/*
chmod -R +x target/$APP_NAME/bin/
cd target && tar -zcvf $APP_NAME.tar.gz $APP_NAME
