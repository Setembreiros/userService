DEV-ENVIRONMENT=development
PROD-ENVIRONMENT=production
DEV-CONN_STR=postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable
PROD-CONN_STR=postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable

update:
	go mod tidy

build: update
	go build -o ./deployment/${PROD-ENVIRONMENT}/userService cmd/main.go

run:
	export CONN_STR="${PROD-CONN_STR}" && export ENVIRONMENT="${PROD-ENVIRONMENT}" && go run cmd/main.go

run-dev:
	export CONN_STR="${DEV-CONN_STR}" && export ENVIRONMENT="${DEV-ENVIRONMENT}" && go run ./cmd/main.go

run-dev-windows: 
	set CONN_STR=${DEV-CONN_STR} && set ENVIRONMENT=${DEV-ENVIRONMENT} && go run ./cmd/main.go

test:
	go generate -v ./internal/... && go test ./internal/...