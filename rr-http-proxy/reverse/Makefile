BIN = $(shell basename $(CURDIR))

run:
	go run main.go

build:
	go build -o $(BIN) main.go

docker:
	docker build .

clean:
	rm -rf $(BIN)

