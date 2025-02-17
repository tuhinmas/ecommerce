GOPATH ?= $(HOME)/go

run-service:
	go run cmd/main.go

format:
	go fmt ./...

run:
	$(GOPATH)/bin/reflex -s -r '\.go$$' make format run-service

test-cov:
	go test -coverprofile=cover.out ./... && go tool cover -html=cover.out -o cover.html

generate-swag:
	swag init -g main.go

generate-jwt-secret:
	$(eval JWT_SECRET := $(shell openssl rand -base64 32))
	@echo "$(JWT_SECRET)"

migration-up:
	migrate -database "mysql://root:opklnm123@tcp(localhost:3307)/contoh" -path migrations up

migration-down:
	migrate -database "mysql://root:opklnm123@tcp(localhost:3307)/contoh" -path migrations down

migration $$(enter):
	@read -p "Migration name:" migration_name; \
	migrate create -ext sql -dir migrations $$migration_name

mock:
	@echo "Generate Mock Interface.."
	mockgen -source="./internal/app/repository/transaction/transaction.go" -destination="./internal/app/repository/transaction/mocks/transaction_mock.go"
	mockgen -source="./internal/app/repository/user/user.go" -destination="./internal/app/repository/user/mocks/user_mock.go"
	mockgen -source="./internal/app/repository/product/product.go" -destination="./internal/app/repository/product/mocks/product_mock.go"
	mockgen -source="./internal/app/repository/warehouse/warehouse.go" -destination="./internal/app/repository/warehouse/mocks/warehouse_mock.go"
	mockgen -source="./internal/app/repository/shop/shop.go" -destination="./internal/app/repository/shop/mocks/shop_mock.go"
	mockgen -source="./internal/app/repository/admin/admin.go" -destination="./internal/app/repository/admin/mocks/admin_mock.go"
	mockgen -source="./internal/app/worker/queue.go" -destination="./internal/app/worker/mocks/queue_mock.go"

lint:
	staticcheck ./...
