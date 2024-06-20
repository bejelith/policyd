GO=go
BUILDDIR=.build
COVERAGEFILE=$(BUILDDIR)/coverage
BIN=$(BUILDDIR)/policyd
TESTCACHE=$(BUILDDIR)/.test_timestmap
PLATFORMS=darwin linux
ARCHS=386 amd64 arm64
MAINPACKAGE=github.com/policyd/pkg/cmd
SOURCES=$(wildcard pkg/**/*.go)

$(BUILDDIR):
	@mkdir -p $(BUILDDIR)

$(TESTCACHE): $(SOURCES) $(BUILDDIR)
	@$(GO) test -coverprofile=$(COVERAGEFILE) ./pkg/...
	@touch $(TESTCACHE)

test: $(TESTCACHE)

coverage: $(COVERAGEFILE)
	$(GO) tool cover -func=$(COVERAGEFILE)
	$(GO) tool cover -html=$(COVERAGEFILE) -o $(BUILDDIR)/coverage.html

$(BIN): $(BUILDDIR)
	@$(GO) build -o $(BIN) $(MAINPACKAGE)
	cp $(BIN) policyd

build: $(BIN) test

build_all: $(BUILDDIR)
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHS), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); $(GO) build -v -o $(BIN)-$(GOOS)-$(GOARCH) $(MAINPACKAGE))))

.PHONY: clean
clean:
	@$(GO) clean -x -i -n -cache -modcache -testcache > /dev/null
	@rm -rf $(BUILDDIR)
