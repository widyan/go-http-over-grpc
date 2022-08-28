#!/bin/bash

echo "Starting generate proto"

go get -u $@ \
  github.com/golang/protobuf/protoc-gen-go \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \

echo "Successfully getting dependency"