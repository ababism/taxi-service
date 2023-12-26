include .env

up:
	docker-compose up --build driver-svc location-svc

up-r:
	docker-compose up driver-svc location-svc

up-d:
	docker-compose up --build driver-svc location-svc -d

down:
	docker-compose down driver-svc location-svc

up-obs:
	docker-compose -f ./deployments/compose.yaml up -d --build

down-obs:
	docker-compose -f ./deployments/compose.yaml down

migrate-up-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_LOCATION_PORT_EXTERNAL)/$(DB_NAME)?sslmode=disable' up

migrate-down-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_LOCATION_PORT_EXTERNAL)/$(DB_NAME)?sslmode=disable' down

migrate-drop-location:
	migrate -path ./projects/location/migrations -database 'postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_LOCATION_PORT_EXTERNAL)/$(DB_NAME)?sslmode=disable' drop

up-kafka:
	docker-compose -f ./kafka/compose.yml up --build -d

down-kafka:
	docker-compose -f ./kafka/compose.yml down