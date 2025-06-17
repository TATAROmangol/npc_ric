installGRPC:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest 

uploadPaths:
	export PATH="$PATH:$(go env GOPATH)/bin"

buildTableProto: installGRPC uploadPaths
	protoc --proto_path=../proto \
	--go_out=./pkg/grpc/table --go_opt=paths=source_relative \
    --go-grpc_out=./pkg/grpc/table --go-grpc_opt=paths=source_relative \
    table.proto