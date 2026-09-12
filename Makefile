BINARY := bin/run-tracker-api
CMD    := ./cmd/main.go

.PHONY: build run test vet fmt tidy clean

build:
	go build -o $(BINARY) $(CMD)

run: build
	$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

clean:
	rm -rf bin
