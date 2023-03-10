
%启动etcd%
start tools/etcd/windows/etcd-v3.5.0-windows-amd64/etcd.exe --config-file tools/etcd/windows/etcd-v3.5.0-windows-amd64/etcd1.conf.yml
start tools/etcd/windows/etcd-v3.5.0-windows-amd64/etcd.exe --config-file tools/etcd/windows/etcd-v3.5.0-windows-amd64/etcd2.conf.yml
start tools/etcd/windows/etcd-v3.5.0-windows-amd64/etcd.exe --config-file tools/etcd/windows/etcd-v3.5.0-windows-amd64/etcd3.conf.yml
%启动网关服%
cd portal
go mod tidy
go build
start portal.exe
cd ..

%启动游戏服%
cd gamesvr
go mod tidy
go build
start gamesvr.exe
cd ..
