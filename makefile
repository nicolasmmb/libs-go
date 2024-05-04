run:
	@go run main.go

build:
	@go build -tags=jsoniter -o bin/main main.go

trace:
	@go tool trace trace.out

bench:
	@echo "Running benchmarks"
	@echo "=================="
	@go test -bench=. ./... -v -count=5 -benchmem

test:
	@echo "Running tests"
	@echo "=================="
	@go test -v ./... -p 6 -v


coverage:
	@go test -coverprofile=./.docs/coverage.out ./...
	@go tool cover -html=./.docs/coverage.out -o ./.docs/coverage.html
	@open ./.docs/coverage.html