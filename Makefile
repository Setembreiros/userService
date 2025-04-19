# existe un ficheiro co mesmo nome do obxectivo ou se considera que o 
# obxectivo está actualizado ao non ter dependencias que cambiaran.
# Neste caso dado que existe un cartafol chamado test no noso proxecto
# o make confundiase e trataba de actualizar este ficheiro en lugares de 
# executar o comando test. Chegaría con ".PHONY: test" neste caso
# pero engado todos por se acaso.
.PHONY: update build run run-dev run-dev-windows test

DEV-ENVIRONMENT=development
PROD-ENVIRONMENT=production
DEV-CONN_STR=postgres://postgres:artis@127.0.0.1:5432/artis?search_path=public&sslmode=disable
PROD-CONN_STR=postgres://postgres:artis12345@artis.c5i8qu2qshhb.eu-west-3.rds.amazonaws.com:5432/artis?search_path=public

update:
	go mod tidy

build: update
	go build -o ./deployment/${PROD-ENVIRONMENT}/userService cmd/main.go

run:
	export CONN_STR="${PROD-CONN_STR}" && export ENVIRONMENT="${PROD-ENVIRONMENT}" && go run ./cmd/main.go

run-dev:
	export CONN_STR="${DEV-CONN_STR}" && export ENVIRONMENT="${DEV-ENVIRONMENT}" && go run ./cmd/main.go

run-dev-windows: 
	set CONN_STR="${DEV-CONN_STR}" && set ENVIRONMENT=${DEV-ENVIRONMENT} && go run ./cmd/main.go

test:
	go generate -v ./internal/... && go test ./internal/...