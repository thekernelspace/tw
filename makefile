BINARY_NAME=tw

install:
	@go build -o ~/.local/bin/$(BINARY_NAME) -v
	@echo "Installed to ~/.local/bin/$(BINARY_NAME)! Type 'tw' at any directory to use it."
