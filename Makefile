NAME      = mattercheck
TARGET    = dist/$(NAME)
DEPS      = $(shell find . -type f -name '*.go')
LDFLAGS	  = -s -w

MANUL_BIN = $(GOPATH)/bin/manul

build     = GOOS=$(1) GOARCH=$(2) go build -ldflags "$(LDFLAGS)" -o $(TARGET)_$(1)_$(2)$(3)

# build a local version for tests, simply run make
mattercheck: $(DEPS)
	go build -o $@

.PHONY: all
all: clean $(DEPS)
	$(call build,darwin,amd64)
	$(call build,freebsd,amd64)
	$(call build,linux,amd64)
	$(call build,linux,arm64)
	$(call build,windows,amd64,.exe)
	cd $(dir $(TARGET)) && sha256sum $(NAME)*
	cd $(dir $(TARGET)) && xz -z9 $(NAME)*

.PHONY: deps
deps: $(MANUL_BIN)
	$(MANUL_BIN) -r -U

.PHONE: clean
clean:
	rm -f mattercheck
	rm -rf dist