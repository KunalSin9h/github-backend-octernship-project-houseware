default:
	@echo "Building and deploying Auth Service"
	cd ./deployments && \
	docker-compose up --build

up:
	@echo "Stoping service if running"
	make down
	@echo "deploying Auth Service"
	cd ./deployments && \
	docker-compose up

down:
	@echo "Stoping Auth Service"
	cd ./deployments && \
	docker-compose down
test:
	@echo "Testing API Endpoints"
	go test ./cmd/api/*.go