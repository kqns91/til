.PHONY: psql-up
psql-up:
		docker compose up -d

.PHONY: psql-exec
psql-exec:
		docker compose exec postgres psql -U postgres

.PHONY: migrate
migrate:
		php artisan migrate

.PHONY: clear
clear:
		php artisan cache:clear
		php artisan config:clear

.PHONY: serve
serve: clear
		php artisan optimize
		php artisan serve

.PHONY: lint
lint:
		vendor/bin/phpcstan --memory-limit=2G

.PHONY: test
test:
		php artisan test
