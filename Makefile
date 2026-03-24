.PHONY: build run test vet clean

BINARY=mdc

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
