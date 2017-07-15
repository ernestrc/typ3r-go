CC=go
CFLAGS=build
TFLAGS=test

REPO=github.com/ernestrc
PKG=typ3r-go
EXEC=typ3r

SRC=$(wildcard **/*.go)
EXECDIR=cmd/typ3r
EXECSRC=$(EXECDIR)/typ3r.go

default: $(EXEC)

.PHONY: clean install

$(EXEC): $(EXECSRC) $(SRC)
	@ find . -name \*.go | xargs cat | grep github.com | awk -F'"' '{print $$2}' | sort | uniq | xargs -L1 go get
	cd $(EXECDIR) && $(CC) $(CFLAGS) -o ../../$(EXEC)

install:
	$(CC) install $(REPO)/$(PKG)/$(EXECDIR)

clean:
	$(CC) clean
	@- rm -f $(EXEC)
