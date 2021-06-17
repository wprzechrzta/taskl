#!/usr/bin/env bash

#echo "Cleaninng up current store"
rm -rf tmp

#CGO_ENABLED=0 GOOS=linux go build -a -v -x -o ./taskl ./cmd/taskl;
#CGO_ENABLED=0 go build  -x -v -o ./taskl ./cmd/taskl;
CGO_ENABLED=0 go build  -o ./taskl ./cmd/taskl;