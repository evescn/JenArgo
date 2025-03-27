#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

docker build -t harbor.evescn.com/devops/jen-argo:v1.$1 -f Dockerfile_$2 .
docker push harbor.evescn.com/devops/jen-argo:v1.$1
