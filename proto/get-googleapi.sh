#!/bin/sh

set -ve

echo "Starting getting googleapi lib"

# apidir="google/api"
# source="https://raw.githubusercontent.com/googleapis/googleapis/master/google/api"

# mkdir -p $apidir

# for file in annotations.proto http.proto; do
#     curl -Ls "${source}/$file" >${apidir}/$file
# done

# protobufdir="lib/google/protobuf"
# source="https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf"

# mkdir -p $protobufdir

# for file in any.proto wrappers.proto; do
#     curl -Ls "${source}/$file" >${protobufdir}/$file
# done

# protobufdir="protoc-gen-openapiv2/options"
# source="https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options"

# mkdir -p $protobufdir

# for file in annotations.proto openapiv2.proto; do
#     curl -Ls "${source}/$file" >${protobufdir}/$file
# done

protobufdir="gogoproto"
source="https://raw.githubusercontent.com/gogo/protobuf/master/gogoproto"

mkdir -p $protobufdir

for file in gogo.proto; do
    curl -Ls "${source}/$file" >${protobufdir}/$file
done

echo "Successfully getting googleapi lib"