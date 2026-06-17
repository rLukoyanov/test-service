.PHONY: build run clean

APP=test-service

build:
	go build -o $(APP) .

run:
	go run .

clean:
	rm -f $(APP)
