.PHONY: dev build test install clean run

dev:
	@./scripts/dev.sh

build:
	@./scripts/build.sh

test:
	@./scripts/test.sh

install:
	@./scripts/install.sh

clean:
	rm -f palettesmith
	rm -rf dist/

run:
	go run ./cmd/palettesmith
