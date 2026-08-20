dev:
	go run main.go

fmt:
	gofumpt -l -w .
	golines -m 80 -w .
	bun fmt

create-schema:
ifndef name
	$(error name is required. Usage: make create-schema name=SchemaName)
endif
	go run entgo.io/ent/cmd/ent new $(name)

generate:
	rm -rf ent/generated && go generate ./ent/...

push:
	go run cmd/push.go

build:
	rm -rf bin && go build -o bin/api
	cd web && bun run build

start:
	./bin/das
