#!/bin/bash
nix-env -iA nixpkgs.openjdk8

sleep 2
# 启动Nacos
cd /home/runner/app/nacos/bin
bash startup.sh -m standalone

sleep 2

# 回到根目录
# cd /

# 启动video-center
cd /home/runner/app/douyin-backend/server/chat-center
go mod tidy
go build ./cmd/main.go
go run ./cmd/main.go &

sleep 4

# 回到根目录
cd /

# 启动chat-center
cd /home/runner/app/douyin-backend/server/video-center
go mod tidy
go build ./cmd/main.go
go run ./cmd/main.go &


sleep 4

# 回到根目录
cd /

# 启动user-center
cd /home/runner/app/douyin-backend/server/user-center
go mod tidy
# cd ./cmd
go build ./cmd/main.go
go run ./cmd/main.go &


sleep 4

# 回到根目录
cd /

# 启动gateway-center
cd /home/runner/app/douyin-backend/server/gateway-center
go mod tidy
go build main.go
go run main.go