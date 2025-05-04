GO=go
BUILDDIR=.build
COVERAGEFILE=$(BUILDDIR)/coverage
ifndef DOCKER
DOCKER=podman
endif
BIN=$(BUILDDIR)/policyd
TESTCACHE=$(BUILDDIR)/.test_timestmap
PLATFORMS=darwin linux
ARCHS=amd64 arm64
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
	@cp $(BIN) policyd

build: test $(BIN)

.PHONY: all
all: build docker $(BUILDDIR)
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHS), $(MAKE) buildarch GOOS=$(GOOS) GOARCH=$(GOARCH);))

buildarch: 
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) build -v -o $(BIN)-$(GOOS)-$(GOARCH) $(MAINPACKAGE)

docker:
	$(foreach GOARCH, $(ARCHS), $(DOCKER) build -f Dockerfile --build-arg BIN=$(BIN)-$(GOARCH) -t policyd-$(GOARCH) .;)

.PHONY: clean
clean:
	@$(GO) clean -x -i -n -cache -modcache -testcache > /dev/null
	@rm -rf $(BUILDDIR)
