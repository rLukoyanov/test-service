.PHONY: build run clean proto

APP=test-service

proto:
	protoc --go_out=. --go_opt=module=test-service \
		--go-grpc_out=. --go-grpc_opt=module=test-service \
		api/proxy.proto

build: proto
	go build -o $(APP) .

run:
	go run .

clean:
	rm -f $(APP)
