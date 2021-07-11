#!/bin/bash

CompileDaemon -build="echo ''" \
    -command="go test ./cmd/... ./internal/... -mod=vendor -v -coverprofile .coverage-unit.txt -tags=unit && go tool cover -func .coverage-unit.txt"


#CompileDaemon -build="echo ''" \
#    -command="go test . -tags=unit"
