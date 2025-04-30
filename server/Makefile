.SILENT:

run:
	go run cmd/talk/main.go --config config/local.env

run-test:
	go run cmd/talk/main.go --config config/test.env

# просто "build" почему-то не работал
build-app:
	go build -o bin/talk -a cmd/talk/main.go && \
	echo "Build completed!"

build-app-with-logs:
	go build -o bin/talk -a -v -x cmd/talk/main.go && \
	echo "Build completed!"

create-migration:
	# Создаст два файла миграции с именем указанным в параметре name или по имени текущей ветки git
	# make create-migration
	# make create-migration name="<migration name>"
	go run cmd/migrator/create/main.go --migration_name=$(or $(name),$(shell git rev-parse --abbrev-ref HEAD))

migrate:
	CONFIG_PATH="config/local.env" go run cmd/migrator/up/main.go

migrate-down-to:
	if [ -z "$(version)" ]; then \
		echo "Error: 'version' parameter is required."; \
		exit 1; \
	fi
	CONFIG_PATH="config/local.env" go run cmd/migrator/down-to/main.go $(version)

# Обновит пакет до последней версии
update-protos:
	go get github.com/Felya-a/chat-app-protos

test:
	WORKDIR_PATH=${shell pwd} CONFIG_PATH=config/test.env \
	ginkgo -v ./...

test-unit:
	WORKDIR_PATH=${shell pwd} CONFIG_PATH=config/local.env \
	ginkgo -v --label-filter=unit ./...

test-unit-watch:
	WORKDIR_PATH=${shell pwd} CONFIG_PATH=config/local.env \
	ginkgo watch -v --label-filter=unit ./...

test-unit-cover:
	WORKDIR_PATH=${shell pwd} CONFIG_PATH=config/local.env \
	ginkgo -v -cover -coverpkg=./... -coverprofile=coverage.out --label-filter=unit ./... && \
	go tool cover -html=coverage.out -o coverage.html && \
	echo "Сформирован файл покрытия - coverage.html"

test-integration:
	WORKDIR_PATH=${shell pwd} CONFIG_PATH=config/test.env \
	ginkgo -v --label-filter=integration ./...

test-watch:
	WORKDIR_PATH=${shell pwd} CONFIG_PATH=config/test.env \
	ginkgo watch -v ./...

test-coverage:
	go test -v -coverprofile=tmp/coverage.txt ./...; \
	go tool cover -html=tmp/coverage.txt -o tmp/coverage.html; \
	echo "Сформирован файл покрытия - tmp/coverage.html"