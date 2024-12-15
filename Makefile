down:
	docker compose down --remove-orphans

run: down
	docker compose up --remove-orphans --build

run_local: down
	docker compose run -d --remove-orphans -p 5432:5432 postgres
	air

test: down
	docker compose -f ./docker-compose.yml -f ./test.docker-compose.yml up --abort-on-container-exit --remove-orphans --build

prod: down
	docker compose -f ./docker-compose.yml -f ./prod.docker-compose.yml up --remove-orphans --build


test_no_docker:
	go test -race ./...
	./fuzz.sh
