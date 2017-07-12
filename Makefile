CC=go
CFLAGS=build
TFLAGS=test
SRC=$(wildcard **/*.go)
TESTDIR=test
EXECDIR=cli
EXEC=typ3r-cli
EXEC_SRC=$(EXECDIR)/main.go

default: $(EXEC)

.PHONY: clean test

$(EXEC): $(EXEC_SRC) $(SRC)
	@ find . -name \*.go | xargs cat | grep github.com | awk -F'"' '{print $$2}' | sort | uniq | xargs -L1 go get
	cd $(EXECDIR) && $(CC) $(CFLAGS) -o ../$(EXEC)

clean:
	@- rm $(EXEC)
