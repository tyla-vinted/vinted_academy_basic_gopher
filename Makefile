SERVER ?=
STEP ?=

.PHONY: run test step online join doctor

## run: open the arena in solo mode (no server)
run:
	go run ./cmd/arena

## test: run all tests
test:
	go test ./avatar/...

## step: run the tests of one step, e.g. make step STEP=3
step:
	@test -n "$(STEP)" || (echo "usage: make step STEP=<1-7>"; exit 1)
	go test ./avatar -run 'Step$(STEP)_' -v

## online: join the workshop lobby without running the tests, e.g. make online SERVER=https://abc.ngrok-free.app
online:
	@test -n "$(SERVER)" || (echo "usage: make online SERVER=https://<host>"; exit 1)
	go run ./cmd/arena -server $(SERVER)

## join: run all tests, then join the workshop lobby
join:
	@test -n "$(SERVER)" || (echo "usage: make join SERVER=https://<host>"; exit 1)
	go vet ./avatar/...
	@go test ./avatar/... || (echo; echo "Your tests are your ticket to the lobby. Make them pass first."; exit 1)
	go run ./cmd/arena -server $(SERVER)

## doctor: check that Go and the build work on this laptop
doctor:
	go version
	go build ./...
	go test ./internal/...
	@echo "All good. Run: make run"
