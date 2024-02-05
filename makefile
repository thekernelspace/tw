BINARY_NAME=tw

.PHONY: build
build:
	@go build -o $(BINARY_NAME)

.PHONY: clean
clean:
	@go clean
	rm -f $(BINARY_NAME)

.PHONY: install
install:
	@go build -o ~/.local/bin/$(BINARY_NAME)
	@echo "Installed to ~/.local/bin/$(BINARY_NAME)! Type 'tw' at any directory to use it."
