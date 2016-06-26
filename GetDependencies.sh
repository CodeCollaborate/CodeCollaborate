#!/usr/bin/env bash

go get github.com/kardianos/govendor
go get gopkg.in/couchbaselabs/gocb.v1
go get github.com/go-sql-driver/mysql
${GOPATH}/bin/govendor sync
