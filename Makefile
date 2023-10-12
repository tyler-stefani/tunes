test-integration:
	docker compose -f docker-compose-test.yml up --build -d \
	&& docker compose -f docker-compose-test.yml logs -f api \
	&& docker compose -f docker-compose-test.yml down