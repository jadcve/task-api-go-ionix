up:
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f api postgres

test:
	go test ./...

run:
	go run ./cmd/api
