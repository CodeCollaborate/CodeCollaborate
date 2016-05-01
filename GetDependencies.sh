#!/usr/bin/env bash

go get github.com/kardianos/govendor
go get gopkg.in/couchbase/gocb.v1
$GOPATH/bin/govendor sync
