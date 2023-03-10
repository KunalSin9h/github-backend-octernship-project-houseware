dev:
	cd deployments && docker-compose down && rm -rf db-data/postgres && docker-compose up -d  && \
	cd ../ && go run cmd/api/*.go