down:
	docker compose down --remove-orphans

run: down
	docker compose up --remove-orphans --build

run_local: down
	docker compose run -d --remove-orphans -p 5432:5432 postgres
	air
