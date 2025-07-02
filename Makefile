
# dependencies:
# go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

test:
go test ./...

generate:
	oapi-codegen -generate types,chi-server -o internal/api/packaging.gen.go -package api openapi.yaml

service-build:
	docker compose build

service-up:
	docker compose up -d

service-down:
	docker compose down
