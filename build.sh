#!/usr/bin/env bash

#CGO_ENABLED=0 GOOS=linux go build -a -v -x -o ./taskl ./cmd/taskl;
CGO_ENABLED=0 go build  -x -v -o ./taskl ./cmd/taskl;