
# dependencies:
# go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

test:
	go test -v -cover -race ./...

# removed for now and left for future
# generate:
# 	oapi-codegen -generate types,chi-server -o internal/api/packaging.gen.go -package api openapi.yaml

service-build:
	docker compose build

service-up:
	docker compose up -d

service-down:
	docker compose down

service-logs:
	docker compose logs -f
