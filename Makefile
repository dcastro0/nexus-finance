.PHONY: up down logs shell

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f app

shell:
	docker compose exec app sh