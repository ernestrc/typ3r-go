CC=go
CFLAGS=build
TFLAGS=test
SRC=$(wildcard *.go)
TESTDIR=test
EXECDIR=bin
EXEC=$(EXECDIR)/typ3r

default: $(EXEC)

$(EXEC): $(EXECDIR)/main.go $(SRC)
	cd $(EXECDIR) && $(CC) $(CFLAGS) -o ../$(EXEC)

.PHONY: clean test

test:
	@ cd $(TESTDIR) && $(CC) $(TFLAGS)

clean:
	@- rm -f $(EXEC)
