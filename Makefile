include .env

up:
	docker-compose up --build driver-svc location-svc

migrate-up-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_LOCATION_PORT_EXTERNAL)/$(DB_NAME)?sslmode=disable' up

migrate-down-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_LOCATION_PORT_EXTERNAL)/$(DB_NAME)?sslmode=disable' down

migrate-drop-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_LOCATION_PORT_EXTERNAL)/$(DB_NAME)?sslmode=disable' drop
