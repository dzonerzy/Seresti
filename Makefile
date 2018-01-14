# Seresti - Makefile
# global variables
GO=$(shell which go)
OUTFILE=serestid
SOURCEDIR=src
INSTALLDIR=/usr/local/bin/


# Do not touch these!
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
DEPS = github.com/go-ini/ini github.com/gorilla/mux

seresti:
	$(info Remember to set GOPATH!)
	$(info Downloading dependencies $(DEPS))
	$(foreach var,$(DEPS),$(GO) get -u $(var);)
	$(GO) build -ldflags="-s -w" -o $(OUTFILE) $(SOURCES)

install:
	install -m 755 $(OUTFILE) $(INSTALLDIR)
	rm -f $(OUTFILE)

.PHONY: clean

clean:
	rm -f $(OUTFILE)

run: seresti
	./$(OUTFILE) --config config/seresti.conf
