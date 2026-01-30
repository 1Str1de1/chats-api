.PHONY: run, build, test, docker-restart
run:
	go run -v ./cmd/app

build:
	go build -v ./cmd/app

docker-restart:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

DEFAULT: build