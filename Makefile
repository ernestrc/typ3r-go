CC=go

ifndef GOBIN
	GOBIN=$(GOPATH)/bin
endif

TARGET=bin
PWD=$(shell pwd)
SRC=$(wildcard **/*.go)
PKGS=.
EXECS=$(sort $(dir $(wildcard cmd/*/)))
EXECSRC=$(wildcard cmd/**/*.go)
EXEC=$(patsubst cmd/%/,$(TARGET)/%,$(EXECS))
GEXEC=$(patsubst cmd/%/,$(GOBIN)/%,$(EXECS))


.PHONY: clean install

default: CHECK $(EXEC)

install: CHECK $(GEXEC) $(PKGS)

clean:
	@-rm -rf $(TARGET)
	@-rm $(GEXEC)

$(TARGET):
	@mkdir $(TARGET)

$(EXEC): $(EXECS) $(EXECSRC) $(SRC) $(TARGET)
	@cd $< && $(CC) build -o $(PWD)/$(patsubst cmd/%,$(TARGET)/%,$@)

$(PKGS): FORCE
	@cd $@ && $(CC) install

$(GEXEC): $(EXECS)
	@cd $< && $(CC) build -o $@

FORCE:

CHECK:
ifndef GOPATH
	$(error GOPATH is undefined)
endif
