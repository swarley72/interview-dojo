.PHONY: proto proto-install clean up up-backend down build image-prune migrate migration superuser seed dump restore

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

migration:
	@read -p "Сервис (core/user): " svc; \
	read -p "Название миграции: " name; \
	migrate create -ext sql -dir $${svc}-service/migrations -seq $$name

superuser:
	DATABASE_URL="postgres://kim:kim@localhost:5432/user_service?sslmode=disable" \
		go run ./user-service/cmd/createsuperuser -login $(login) -password $(password)

seed:
	docker compose exec -T postgres psql -U kim -d core_service < core-service/seeds/seed_questions.sql

down:
	docker compose down

build:
	docker compose build

image-prune:
	docker image prune

dump:
	mkdir -p dumps
	docker compose exec -T postgres pg_dump -U kim --clean core_service | zstd -o dumps/core_service_$(shell date +%Y%m%d_%H%M%S).sql.zst

restore:
	zstd -dc $(or $(file),$(shell ls -t dumps/*.sql.zst | head -1)) | docker compose exec -T postgres psql -U kim -d core_service
