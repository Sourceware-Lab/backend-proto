run:
	docker compose down --remove-orphans
	docker compose up --remove-orphans --build

run_local:
	docker compose down--remove-orphans
	docker compose run -d --remove-orphans -p 5432:5432 postgres
	air
