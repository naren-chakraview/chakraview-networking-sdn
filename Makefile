.PHONY: help build test clean lab-up lab-down

help:
	@echo "Targets:"
	@echo "  make build          Build all components (foundation, controller, fabric)"
	@echo "  make test           Run all tests"
	@echo "  make clean          Clean build artifacts"
	@echo "  make lab-up         Start Docker Compose lab"
	@echo "  make lab-down       Stop Docker Compose lab"
	@echo "  make lab-init       Initialize lab (requires lab-up)"
	@echo "  make docs           Build MkDocs site"

build:
	cd foundation && make build
	cd controller && go build -o sdn-controller ./cmd/sdn-controller
	cd fabric && pip install -r requirements.txt

test:
	cd foundation && make test
	cd controller && go test ./...
	cd fabric && pytest

clean:
	cd foundation && make clean
	rm -f controller/sdn-controller
	cd fabric && find . -name __pycache__ -exec rm -rf {} + 2>/dev/null || true

lab-up:
	cd lab && docker-compose up -d

lab-down:
	cd lab && docker-compose down

lab-init:
	cd lab && bash scripts/init.sh

docs:
	pip install mkdocs mkdocs-material
	mkdocs build
