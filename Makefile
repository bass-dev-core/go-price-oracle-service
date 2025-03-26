APP_NAME=go-price-oracle-service
PORT=6942

OS_LIST=linux darwin windows
ARCH_LIST=amd64

ifeq ($(OS),Windows_NT)
	EXE_EXT=.exe
	SHELL := cmd
	.SHELLFLAGS := /C
	NULL := NUL
else
	EXE_EXT=
	SHELL := /bin/sh
	.SHELLFLAGS := -c
	NULL := /dev/null
endif

install:
	go clean -modcache
	go mod tidy
	go mod vendor	
	npm i -g typescript
	npm i

prepare: install build

build:
	go build -o $(APP_NAME)$(EXE_EXT) ./cmd

run:
	go run ./cmd

test:
	go test ./...

build-all:
	@echo "🔧 Starting cross-build..."
	@for OS in $(OS_LIST); do \
		for ARCH in $(ARCH_LIST); do \
			OUT_DIR=build/$$OS-$$ARCH; \
			mkdir -p $$OUT_DIR; \
			EXT=$$( [ "$$OS" = "windows" ] && echo ".exe" || echo "" ); \
			GOOS=$$OS GOARCH=$$ARCH go build -mod=vendor -o $$OUT_DIR/$(APP_NAME)$$EXT ./cmd; \
			echo "✅ Built: $$OUT_DIR/$(APP_NAME)$$EXT"; \
		done; \
	done

clean:
	rm -rf build vendor $(APP_NAME) $(APP_NAME).exe

deploy: docker-build docker-run

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -p $(PORT):6942 $(APP_NAME)
