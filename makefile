serve:
	go run main.go	
generate-mocks:
	go generate ./...	
test:
	go test ./...	
make docker-dev:
	docker compose up -d couchserver
docker-run:
	docker compose up -d