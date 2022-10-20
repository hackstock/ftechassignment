BINARY_NAME=barbershop

build: deps
	go build -o ${BINARY_NAME} -ldflags="-s -w" cmd/barbershop/main.go

run: deps build
	./${BINARY_NAME}

test:
	go test -v -count=1 -coverprofile=coverage.out ./...

clean:
	go clean
	rm ${BINARY_NAME}

deps:
	go mod download