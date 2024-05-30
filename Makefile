DEV-ENVIRONMENT=development
PROD-ENVIRONMENT=production

update:
	go mod tidy
	
build:
    go build -o ./deplotment/${PROD-ENVIRONMENT}/nome-do-teu-proxecto ./cmd

run:
    go run ./cmd/main.go

test:
    go test ./...