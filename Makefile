ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
APP=arche-api
APP_VERSION:="1.0"
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

build: fmt vet
	go build -o $(APP_EXECUTABLE) main/*.go

serve: fmt vet
	go run main/*.go
