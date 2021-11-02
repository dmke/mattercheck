NAME      = mattercheck
TARGET    = dist/$(NAME)
DEPS      = $(shell find . -type f -name '*.go')
LDFLAGS	  = -s -w

export GO111MODULE=on
build     = GOOS=$(1) GOARCH=$(2) go build -ldflags "$(LDFLAGS)" -trimpath -o $(TARGET)_$(1)_$(2)$(3)

# build a local version for tests, simply run make
mattercheck: $(DEPS)
	go build -ldflags "$(LDFLAGS)" -o $@

.PHONY: all
all: clean $(DEPS)
	$(call build,darwin,amd64)
	$(call build,freebsd,amd64)
	$(call build,linux,amd64)
	$(call build,linux,arm64)
	$(call build,windows,amd64,.exe)
	cd $(dir $(TARGET)) && sha256sum $(NAME)*
	cd $(dir $(TARGET)) && xz -z9 $(NAME)*

.PHONY: install
install: $(DEPS)
	go install -ldflags "$(LDFLAGS)" .

.PHONE: clean
clean:
	rm -f mattercheck
	rm -rf dist

.PHONY: update-fixture
update-fixture:
	curl -sSLo releases/testdata/version-archive.html \
		https://docs.mattermost.com/upgrade/version-archive.html
