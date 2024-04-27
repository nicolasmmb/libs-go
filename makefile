run:
	@go run main.go

build:
	@go build -tags=jsoniter -o bin/main main.go

trace:
	@go tool trace trace.out

bench:
	@echo "Running benchmarks"
	@echo "=================="
	@go test -bench=. ./repository -v -count=5 -benchmem
	@echo "=================="
	@go test -bench=. ./bus -v -count=5 -benchmem
	@echo "=================="
	@go test -bench=. ./uow -v -count=5 -benchmem
	@echo "=================="

test:
	@echo "Running tests"
	@echo "=================="
	@go test -v ./repository
	@echo "=================="
	@go test -v ./bus
	@echo "=================="
	@go test -v ./uow
	@echo "=================="
