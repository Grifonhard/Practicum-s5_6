.PHONY: build
build:
	cd cmd/accrual && go build -o accrual main.go
	cd cmd/gophermart && go build -o gophermart main.go

.PHONY: tidy
tidy:
	go mod tidy
