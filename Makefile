
# e2e 测试
orm_e2e:
	docker compose down
	docker compose up -d
	go test -race -tags=e2e ./internal/integration/
	docker compose down