.PHONY: build up down swagger redis-cli redis-keys logs help


build:
	go build -o bin/server cmd/main.go

up:
	docker-compose up --build 

down:
	docker-compose down

swagger:
	swag init -g cmd/main.go


redis-keys:
	docker-compose exec redis redis-cli KEYS "sent_message:*"

logs:
	docker-compose logs -f app

psql:
	docker-compose exec postgres psql -U postgres -d insider_messaging