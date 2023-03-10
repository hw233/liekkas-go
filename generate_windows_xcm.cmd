

%导表--------------------%
set ConfigPath= C:\code\overlord-config-backup

cd tools/excel_to_csv/
go mod tidy
go run main.go e2c.go -input %ConfigPath% -output ../../shared/csv/data -constOutput ../../shared/csv/static -constTpl const.tpl
cd ../../

cd tools/csv2go/
go mod tidy
cd ../../
go run tools/csv2go/convert.go tools/csv2go/generator.go tools/csv2go/main.go tools/csv2go/manager_generator.go tools/csv2go/parser.go ./shared/csv/base ./shared/csv/data

cd ./shared/csv/base && go fmt
cd ../../../
cd ./shared/csv/static && go fmt
cd ../../../


set ProtoPath=./shared/protobuf/proto/
for /f "delims=\" %%a in ('dir /b /a-d /o-d "%ProtoPath%\*.*"') do (
  start ./tools/buf/windows/protoc.exe --plugin=protoc-gen-go=./tools/protobuf/windows/protoc-gen-go.exe --plugin=protoc-gen-go-grpc=./tools/protobuf/windows/protoc-gen-go-grpc.exe  --go-grpc_out=./shared/protobuf/pb --go_out=./shared/protobuf/pb ./shared/protobuf/proto/%%a

)
pause
