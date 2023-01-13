up:
	docker-compose up

down:
	docker-compose down

build:
	go build -ldflags="-w -s"

run:
	go run main.go
