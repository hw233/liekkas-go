# !/bin/bash
# Auth: 雪辙<pengqing@bilibli.com>
# Release: v1.0 2021-07-23
# Desc: 生成配置数据
# Params: None
# Usage: ./generate_macOS.sh

# 生成csv
./tools/excel_to_csv/macOS/excel_to_csv -input=/Users/Sora/WorkSpace/overlord/overlord-config -output=./shared/csv/data -constOutput=./shared/csv/static -constTpl=./tools/excel_to_csv/const.tpl
./tools/csv2go/macOS/csv2go ./shared/csv/base ./shared/csv/data
cd ./shared/csv/base && go fmt
cd -
cd ./shared/csv/static && go fmt
cd -
#gofmt -l -s ./shared/csv/base ./shared/csv/static

# 生成protobuf
./tools/buf/macOS/buf generate --path ./shared/protobuf/proto --template '{"version":"v1beta1","plugins":[{"name":"go","path":"./tools/protobuf/macOS/protoc-gen-go","out":"./shared/protobuf/pb"},{"name":"go-grpc","path":"./tools/protobuf/macOS/protoc-gen-go-grpc","out":"./shared/protobuf/pb"}]}'

# 更新客户端表格
cp ./shared/csv/data/protocol.csv /Users/Sora/WorkSpace/overlord/overlord-protocol/proto/
cp ./shared/protobuf/proto/*.proto /Users/Sora/WorkSpace/overlord/overlord-protocol/proto/
/Users/Sora/WorkSpace/overlord/overlord-protocol/ProtocolTooMac.app/Contents/MacOS/ProtocolTool
cd /Users/Sora/WorkSpace/overlord/overlord-protocol/
java -jar ./ProtoGenerator.jar
cd -