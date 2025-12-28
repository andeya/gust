
BIN=go
# Limit compiler concurrency to reduce memory usage (default: 4)
# Set COVERAGE_MAXPROCS environment variable to override
COVERAGE_MAXPROCS ?= 4
# Limit parallel test execution (default: 2 to balance speed and memory)
# Set COVERAGE_PARALLEL environment variable to override
COVERAGE_PARALLEL ?= 2

.PHONY: all build test lint coverage coverage-modular clean-coverage

all: build test lint

build:
	${BIN} build -v $$(${BIN} list ./... | grep -v '/scripts/')

test:
	go test -v $$(${BIN} list ./... | grep -v '/scripts/')
watch-test:
	reflex -t 50ms -s -- sh -c 'gotest -v $$(go list ./... | grep -v "/scripts/")'

bench:
	go test -benchmem -count 3 -bench $$(${BIN} list ./... | grep -v '/scripts/')
watch-bench:
	reflex -t 50ms -s -- sh -c 'go test -benchmem -count 3 -bench $$(${BIN} list ./... | grep -v "/scripts/")'

# Legacy coverage (runs all packages at once, high memory usage)
coverage:
	${BIN} test -v -coverprofile=cover.out -covermode=atomic $$(${BIN} list ./... | grep -v '/scripts/')
	${BIN} tool cover -html=cover.out -o cover.html

# Modular coverage (runs each package separately to reduce memory usage)
# Limits compiler concurrency to reduce memory usage during compilation
coverage-modular:
	@echo "Starting modular coverage testing..."
	@mkdir -p coverage-tmp
	@rm -f coverage-tmp/*.out
	@echo "Getting list of all packages (excluding scripts)..."
	@${BIN} list ./... | grep -v '/scripts/' | grep -v 'scripts' | grep -v "^github.com/andeya/gust$$" > coverage-tmp/packages.txt || true
	@echo "Packages to test:"
	@cat coverage-tmp/packages.txt
	@echo "Running tests for each package separately (GOMAXPROCS=$(COVERAGE_MAXPROCS), parallel=$(COVERAGE_PARALLEL))..."
	@while IFS= read -r pkg; do \
		if [ -n "$$pkg" ]; then \
			pkg_name=$$(echo $$pkg | sed 's|github.com/andeya/gust/||' | sed 's|/|_|g' | sed 's|^$$|root|'); \
			echo "Testing package: $$pkg (output: coverage-tmp/$$pkg_name.out)"; \
			GOMAXPROCS=$(COVERAGE_MAXPROCS) ${BIN} test -p $(COVERAGE_PARALLEL) -v -coverprofile=coverage-tmp/$$pkg_name.out -covermode=atomic $$pkg || exit 1; \
		fi; \
	done < coverage-tmp/packages.txt
	@echo "Merging all coverage files..."
	@go run scripts/merge_coverage.go cover.out coverage-tmp/*.out || (echo "Error: Failed to merge coverage files" && exit 1)
	@echo "Generating HTML coverage report..."
	@${BIN} tool cover -html=cover.out -o cover.html
	@echo ""
	@echo "=========================================="
	@echo "Overall Coverage Summary:"
	@echo "=========================================="
	@${BIN} tool cover -func=cover.out | tail -1
	@echo "=========================================="
	@echo ""
	@echo "Coverage testing completed! Results saved in cover.out and cover.html"
	@echo "Cleaning up temporary files..."
	@rm -rf coverage-tmp

# Clean coverage temporary files
clean-coverage:
	rm -rf coverage-tmp
	rm -f cover.out cover.html *.out

# tools
tools:
	${BIN} install github.com/cespare/reflex@latest
	${BIN} install github.com/rakyll/gotest@latest
	${BIN} install github.com/psampaz/go-mod-outdated@latest
	${BIN} install github.com/jondot/goweight@latest
	${BIN} install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	${BIN} get -t -u golang.org/x/tools/cmd/cover
	${BIN} install github.com/sonatype-nexus-community/nancy@latest
	go mod tidy

# Build coverage merge tool (optional, for faster execution)
build-coverage-merge-tool:
	${BIN} build -o scripts/merge_coverage scripts/merge_coverage.go

lint:
	golangci-lint run --print-resources-usage --fast --timeout 20m0s --max-same-issues 50 $$(${BIN} list ./... | grep -v '/scripts/')
lint-fix:
	golangci-lint run --print-resources-usage --fast --timeout 20m0s --max-same-issues 50 --fix $$(${BIN} list ./... | grep -v '/scripts/')

audit:
	${BIN} mod tidy
	${BIN} list -json -m all | nancy sleuth

outdated:
	${BIN} mod tidy
	${BIN} list -u -m -json all | go-mod-outdated -update -direct

weight:
	goweight