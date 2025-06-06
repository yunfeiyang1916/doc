package v1

//go:generate protoc --proto_path=../../../third_party -I . --go_out=. --go-grpc_out=. *.proto
