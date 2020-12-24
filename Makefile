GOCMD=go
GOBUILD=$(GOCMD) build

image: build run
page: build run-page

build:
	@echo "Building..."
	@$(GOBUILD)

run:
	@echo "Downloading images..."
	@./kuro

run-page:
	@echo "Downloading pages..."
	@./kuro -page