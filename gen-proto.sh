#!/bin/bash

MODULE_OPT="--go_opt=module=github.com/the-tinderbox/tron-grpc --go-grpc_opt=module=github.com/the-tinderbox/tron-grpc"

rm -rf core api

pushd proto

protoc -I. --go_out=../ --go-grpc_out=../ $MODULE_OPT api/*.proto core/*.proto core/contract/*.proto

popd
