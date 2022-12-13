#!/bin/bash

RUN_NAME="xiangqin.user.api"

export GO111MODULE=on
go mod download
mkdir -p output/bin output/conf
cp script/bootstrap.sh script/pre_nginx.sh script/settings.py output 2>/dev/null
chmod +x output/bootstrap.sh output/pre_nginx.sh
cp script/bootstrap.sh output/bootstrap_staging.sh
chmod +x output/bootstrap_staging.sh
find output/conf -type f ! -name "*_local.*" | xargs -I{} cp {} output/conf/

go build -o output/bin/${RUN_NAME}