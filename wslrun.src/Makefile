.PHONY: install build

build:
	@go build -o "wslrun" -ldflags "-s -w"

install:
	@cp -f "wslrun" /usr/local/bin
	@cp -f "wslrun-open" /usr/local/bin

