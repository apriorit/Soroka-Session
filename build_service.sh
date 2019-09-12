#!/usr/bin/env bash

## run go test
echo Running go test...
go test ./...

CGO_ENABLED=0 GOOS=linux go build ${GOFLAGS} -a \
    -installsuffix cgo \
    github.com/Soroka-EDMS/svc/sessions/cmd/sessionssvc/