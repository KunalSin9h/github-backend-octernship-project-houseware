# Build and deploy Auth Service
# This will start both auth service and postgres database
default:
	@echo "Building and deploying Auth Service"
	cd ./deployments && \
	docker-compose up --build

# deploy auth service
# this will start both auth service and postgres database without building
up:
	@echo "Stoping service if running"
	make down
	@echo "deploying Auth Service"
	cd ./deployments && \
	docker-compose up

# stop auth service
# this will stop both auth service and postgres database
down:
	@echo "Stoping Auth Service"
	cd ./deployments && \
	docker-compose down

# build auth service binary
# this will build auth service binary
build:
	@echo "Building Auth Service Binary"
	go build -o main ./cmd/api/*.go

# run tests
# this will run all the tests
test:
	@echo "Testing API Endpoints"
	go test ./cmd/api/*.go