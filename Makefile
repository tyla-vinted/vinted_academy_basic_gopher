SERVER ?=

.PHONY: run check next test online join doctor

## run: open the arena in solo mode (no server)
run:
	go run ./cmd

## check: test the steps you have
check:
	@go run ./internal/steps check

## next: if your steps pass, unlock the next step
next:
	@go run ./internal/steps next

## test: run all tests, with less text than make check
test:
	go test ./avatar/...

## online: join the workshop lobby, e.g. make online SERVER=https://abc.ngrok-free.app
online:
	@test -n "$(SERVER)" || (echo "usage: make online SERVER=https://<host>"; exit 1)
	go run ./cmd -server $(SERVER)

## join: when all 7 steps pass, join the workshop lobby for the finale
join:
	@test -n "$(SERVER)" || (echo "usage: make join SERVER=https://<host>"; exit 1)
	@go run ./internal/steps final
	go run ./cmd -server $(SERVER)

## doctor: check that Go and the build work on this laptop
doctor:
	go version
	go build ./...
	go test ./internal/...
	@echo "All good. Run: make run"
