package pb

//go:generate protoc --go_out=. --go-grpc_out=. *.proto

// 下面的命令需要拷贝到命令行单独执行，否则会报错
// protoc -I . --go_out=":." --validate_out="lang=go:." *.proto
