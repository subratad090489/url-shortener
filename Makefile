.PHONY: build run test test-verbose test-coverage test-shortener docker-build docker-run docker-compose-up docker-compose-down clean
# Go commands
build:
	go build -o bin/url-shortener ./cmd/server


run:
	go run ./cmd/server

# Test commands
test:
	go test ./...

test-verbose:
	go test -v ./...


# Docker commands
docker-build:
	docker build -t url-shortener:latest -f deployments/docker/Dockerfile .


docker-run:
	docker run -p 8080:8080 --name url-shortener url-shortener:latest


# Docker Compose commands
docker-compose-up:
	docker-compose -f deployments/docker/docker-compose.yml up --build

docker-compose-down:
	docker-compose -f deployments/docker/docker-compose.yml down

# Clean up
clean:
	rm -rf bin/
	docker-compose -f deployments/docker/docker-compose.yml down
	docker rm -f url-shortener 2>/dev/null || true
	docker rmi url-shortener:latest 2>/dev/null || true
