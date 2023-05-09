# gRPC

## Goコードの生成

```cmd
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

```cmd
protoc --proto_path=. --go_out=. --go-grpc_out=. ./*.proto
```

## Client Streaming RPC
- クライアントから複数回リクエストを送信し、サーバーがそれに対してレスポンスを1回返す通信方式
- クライアントからサーバーに向けて大きいサイズのファイルをアップロードする場合などに使用する

.protoのrpcのリクエストに`stream`をつけるだけ

```proto
service ExampleService {
	rpc ClientStream (stream ExampleRequest) returns (ExampleResponse);
}
```

- クライアント
  - streamclientを生成
  - `client.Send(req)`で1回分の送信
  - 送り切ったら、`res, err := client.CloseAndRecv()`でレスポンスを受け取ってコネクションを閉じる

- サーバー
  - `stream.Recv()`で1回分のデータを受信
  - 受け取り切ったら`stream.SendAndClose()`でレスポンスを返しコネクションを閉じる