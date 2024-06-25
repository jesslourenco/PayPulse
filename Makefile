build:
	go build -o ./bin/gopay cmd/gopay/main.go

create-mocks:
	mockery