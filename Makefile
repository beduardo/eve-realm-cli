GIT_HASH   := $(shell git rev-parse --short HEAD)
VERSION    := $(shell cat VERSION)
BUILD_DATE := $(shell date -u +%Y-%m-%d)

BINARY_NAME := eve-realm
CMD_PATH    := ./cmd
DIST_DIR    := dist
INSTALL_DIR := /usr/local/bin

LDFLAGS := -X main.Version=$(VERSION) -X main.GitHash=$(GIT_HASH) -X main.BuildDate=$(BUILD_DATE)

.PHONY: build test install clean version \
        bump-patch bump-minor bump-major \
        release-patch release-minor release-major

# ─── Build ───────────────────────────────────────────────────────────────────

build:
	mkdir -p $(DIST_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME) $(CMD_PATH)

test:
	go test -count=1 ./...

install:
	@test -f $(DIST_DIR)/$(BINARY_NAME) || { echo "No binary in $(DIST_DIR)/. Run 'make build' first."; exit 1; }
	@test -w $(INSTALL_DIR) || { echo "$(INSTALL_DIR) is not writable. Run 'sudo make install' or override INSTALL_DIR."; exit 1; }
	cp $(DIST_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(DIST_DIR)

version:
	@echo $(VERSION)

# ─── Version Bumping ─────────────────────────────────────────────────────────

bump-patch:
	@version=$$(cat VERSION); \
	major=$$(echo $$version | awk -F. '{print $$1}'); \
	minor=$$(echo $$version | awk -F. '{print $$2}'); \
	patch=$$(echo $$version | awk -F. '{print $$3}'); \
	patch=$$((patch + 1)); \
	echo "$$major.$$minor.$$patch" > VERSION; \
	echo "Version bumped to $$(cat VERSION)"

bump-minor:
	@version=$$(cat VERSION); \
	major=$$(echo $$version | awk -F. '{print $$1}'); \
	minor=$$(echo $$version | awk -F. '{print $$2}'); \
	minor=$$((minor + 1)); \
	echo "$$major.$$minor.0" > VERSION; \
	echo "Version bumped to $$(cat VERSION)"

bump-major:
	@version=$$(cat VERSION); \
	major=$$(echo $$version | awk -F. '{print $$1}'); \
	major=$$((major + 1)); \
	echo "$$major.0.0" > VERSION; \
	echo "Version bumped to $$(cat VERSION)"

# ─── Release ─────────────────────────────────────────────────────────────────

release-patch: test bump-patch
	@$(MAKE) build
	@$(MAKE) install
	@echo "Released $$(cat VERSION) to $(INSTALL_DIR)/$(BINARY_NAME)"
	@$(INSTALL_DIR)/$(BINARY_NAME) version

release-minor: test bump-minor
	@$(MAKE) build
	@$(MAKE) install
	@echo "Released $$(cat VERSION) to $(INSTALL_DIR)/$(BINARY_NAME)"
	@$(INSTALL_DIR)/$(BINARY_NAME) version

release-major: test bump-major
	@$(MAKE) build
	@$(MAKE) install
	@echo "Released $$(cat VERSION) to $(INSTALL_DIR)/$(BINARY_NAME)"
	@$(INSTALL_DIR)/$(BINARY_NAME) version
