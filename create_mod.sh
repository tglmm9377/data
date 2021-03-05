#! /bin/bash
mkdir api && cd api
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod init api
go get -u github.com/gin-gonic/gin

