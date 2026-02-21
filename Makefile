.PHONY: proto proto-install clean up up-backend down build image-prune migrate superuser

proto:
	protoc --proto_path=. --proto_path=proto \
		   --go_out=. --go_opt=paths=source_relative \
		   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		   proto/user/user.proto proto/core/core.proto

clean:
	find proto -name "*.pb.go" -delete

proto-install:
	cd proto \
	&& go install google.golang.org/protobuf/cmd/protoc-gen-go \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

up:
	docker compose up --build -d

up-backend:
	docker compose up --build -d postgres user-service core-service api-gateway

migrate:
	migrate -path user-service/migrations -database "postgres://kim:kim@localhost:5432/user_service?sslmode=disable" up
	migrate -path core-service/migrations -database "postgres://kim:kim@localhost:5432/core_service?sslmode=disable" up

superuser:
	DATABASE_URL="postgres://kim:kim@localhost:5432/user_service?sslmode=disable" \
		go run ./user-service/cmd/createsuperuser -login $(login) -password $(password)

down:
	docker compose down

build:
	docker compose build

image-prune:
	docker image prune
