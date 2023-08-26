# Makefile

.PHONY: lint
lint:
	golangci-lint run ${args} ./ ...

.PHONY: lint-fix
lint-fix:
	@make lint args=' --fix -v' cons_args='-v'

.PHONY: build
build:
	@./scripts/build.sh

.PHONY: run
run:
	go run src/cmd/server/main.go