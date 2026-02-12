.PHONY: proto proto-intall clean

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		   proto/user/user.proto proto/core/core.proto

clean:
	find proto -name "*.pb.go" -delete

proto-install:
	cd proto \
	&& go install google.golang.org/protobuf/cmd/protoc-gen-go \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
