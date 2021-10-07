GO=go
GOFLAGS=

.PHONY: dist/user dist/auth dist/tb
all: dist/user dist/auth dist/tb

dist/user: user/main.go
	$(GO) $(GOFLAGS) build -o $@ user/main.go

dist/auth: auth/main.go
	$(GO) $(GOFLAGS) build -o $@ auth/main.go

dist/tb: tb/main.go
	$(GO) $(GOFLAGS) build -o $@ tb/main.go

test: 
	$(GO) test ./user ./auth ./tb
