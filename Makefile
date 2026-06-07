run-sample:
	go run ./cmd/ingestor/main.go ./data/1.csv

docker-build:
	docker build -t lei_explorer:latest -f infra/Dockerfile .

# Run the application using Docker and pass csv file from host to container
docker-run:
	docker run -p 8080:8080 \
		-v $$PWD/data:/root/data:ro lei_explorer:latest \
		./app ./data/1.csv

up:
	docker-compose -f infra/docker-compose.yaml up --build

down:
	docker-compose -f infra/docker-compose.yaml down -v