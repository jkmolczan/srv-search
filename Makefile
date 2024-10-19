# Makefile

run:
	go run app/main.go app/config.go

test:
	go test ./pkg/...
