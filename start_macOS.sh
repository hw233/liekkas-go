# !/bin/bash
# Auth: 雪辙<pengqing@bilibli.com>
# Release: v1.0 2021-07-23
# Desc: 启动服务器
# Params: None
# Usage: ./start_macOS.sh

# 启动etcd
nohup ./tools/etcd/macOS/etcd-v3.5.0-darwin-amd64/etcd --config-file=./tools/etcd/macOS/etcd-v3.5.0-darwin-amd64/etcd1.conf.yml &
nohup ./tools/etcd/macOS/etcd-v3.5.0-darwin-amd64/etcd --config-file=./tools/etcd/macOS/etcd-v3.5.0-darwin-amd64/etcd2.conf.yml &
nohup ./tools/etcd/macOS/etcd-v3.5.0-darwin-amd64/etcd --config-file=./tools/etcd/macOS/etcd-v3.5.0-darwin-amd64/etcd3.conf.yml &