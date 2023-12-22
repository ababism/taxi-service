include .env

up:
	docker-compose up --build driver-svc location-svc

migrate-up-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(DB_LOCATION_PORT_EXTERNAL)/postgres?sslmode=disable' up

migrate-down-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(DB_LOCATION_PORT_EXTERNAL)/postgres?sslmode=disable' down

migrate-drop-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(DB_LOCATION_PORT_EXTERNAL)/postgres?sslmode=disable' drop
