# !/bin/bash
# Auth: 雪辙<pengqing@bilibli.com>
# Release: v1.0 2021-07-23
# Desc: 生成配置数据
# Params: None
# Usage: ./gen_proto_macOS.sh

../tools/buf/macOS/buf generate --path ./pure.proto --template '{"version":"v1beta1","plugins":[{"name":"go","path":"../tools/protobuf/macOS/protoc-gen-go","out":"./pb"},{"name":"go-grpc","path":"../tools/protobuf/macOS/protoc-gen-go-grpc","out":"./pb"}]}'