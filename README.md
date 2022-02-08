# gRPC Demo

## 安裝

### Proto Buffer CLI 工具
```shell
$ brew install protobuf

# 檢查 protobuf 版本, 確認版本為 3+
$ protoc --version
```

## 編譯 & 建置

### .proto 檔案 -> .go 檔案
```shell
# 將目錄切換到有 .proto 檔案的目錄下
$ cd proto

# 編譯
$ protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative *.proto
```

### server application
```shell
$ go build -o server server.go
```

### client application
```shell
$ go build -o client client.go
```

## 使用

### 啟動 gRPC Server
```shell
$ ./server
```

### 啟動 gRPC Client

提供四種 gRPC 服務類型
- Unary
- Server Stream
- Client Stream
- Bi-direction Stream

```shell
Usage: ./client [options]
Options:
  -type string
    	Server type [unary, server, client, bi-direct] (default "unary")
  -value int
    	Value for request (default 5)
```