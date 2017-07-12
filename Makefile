CC=go
CFLAGS=build
TFLAGS=test
SRC=$(wildcard **/*.go)
TESTDIR=test
EXECDIR=cli
EXEC=typ3r-cli

default: $(EXEC)

$(EXEC): $(EXECDIR)/main.go $(SRC)
	cd $(EXECDIR) && $(CC) $(CFLAGS) -o ../$(EXEC)

.PHONY: clean test

test:
	@ cd $(TESTDIR) && $(CC) $(TFLAGS)

clean:
	@- rm $(EXEC)
