.PHONY: build
build:
	go build -o bin/dac main.go

.PHONY: clean
clean:
	rm -rf bin
