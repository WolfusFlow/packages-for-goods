
# dependencies:
# go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

test:
	go test -v -cover -race ./...

# removed for now and left for future
# generate:
# 	oapi-codegen -generate types,chi-server -o internal/api/packaging.gen.go -package api openapi.yaml

service-build:
	docker compose -f docker-compose.yaml  build

service-up:
	docker compose -f docker-compose.yaml  up -d

service-down:
	docker compose -f docker-compose.yaml  down

service-logs:
	docker compose -f docker-compose.yaml logs -f

atlas-init:
	atlas migrate diff initial \
  	  --to file://infra/atlas/schema.sql \
  	  --dev-url "docker://postgres/15/dev?search_path=public" \
  	  --dir file://infra/atlas/migrations \
  	  --format '{{ sql . "  " }}'

atlas-diff:
	atlas migrate diff $(name) \
		--to file://infra/atlas/schema.sql \
		--dev-url "docker://postgres/15/dev?search_path=public" \
		--dir file://infra/atlas/migrations \
		--format '{{ sql . "  " }}'

seed:
	docker compose run --rm seed
