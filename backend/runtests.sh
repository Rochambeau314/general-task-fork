#!/bin/bash

export DB_NAME=test
go test -v ./...
[ $? -eq 0 ]  || exit 1  # exits if test failures
unset DB_NAME
