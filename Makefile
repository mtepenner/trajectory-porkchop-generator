.PHONY: up down build test-lambert fmt

up:
	docker-compose up --build

down:
	docker-compose down

build:
	cd compute_engine && go build ./...
	cd api_gateway && pip install -r requirements.txt
	cd mission_planner_ui && npm install && npm run build

test-lambert:
	cd compute_engine && go test ./internal/orbital_math/...

fmt:
	cd compute_engine && gofmt -w .
	cd api_gateway && python -m black app/

clean:
	docker-compose down -v
	rm -rf compute_engine/solver mission_planner_ui/build
