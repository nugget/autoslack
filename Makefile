BINDIR		?=$(HOME)/bin
CONFDIR		?=$(HOME)/.config/autoslack

autoslack: go.mod go.sum main.go
	go build .

install: autoslack config.json
	@echo "Installing autoslack and user config"
	@install -v -d -m 0750 $(CONFDIR)
	@install -v -d -m 0755 $(BINDIR)
	@install -v -C -p -m 0644 config.json $(CONFDIR)
	@install -v -C -p -m 0755 autoslack_launch $(BINDIR)
	@install -v -C -p -m 0755 autoslack $(BINDIR)
