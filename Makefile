WORKDAY_COUNTER_BIN := bin/workday-counter

GO ?= go
GO_MD2MAN ?= go-md2man

VERSION := $(shell cat VERSION)
USE_VENDOR =
LOCAL_LDFLAGS = -buildmode=pie -ldflags "-X=main.version=$(VERSION)"

.PHONY: all api build vendor
all: dep build

dep: ## Get the dependencies
	@$(GO) get -v -d ./...

update: ## Get and update the dependencies
	@$(GO) get -v -d -u ./...

tidy: ## Clean up dependencies
	@$(GO) mod tidy

vendor: dep ## Create vendor directory
	@$(GO) mod vendor

build: ## Build the binary files
	$(GO) build -v -o $(WORKDAY_COUNTER_BIN) $(USE_VENDOR) $(LOCAL_LDFLAGS) ./cmd/workday-counter

clean: ## Remove previous builds
	@rm -f $(WORKDAY_COUNTER_BIN)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: release
release: ## create release package from git
	git clone https://github.com/thkukuk/workday-counter
	mv workday-counter workday-counter-$(VERSION)
	#sed -i -e 's|USE_VENDOR =|USE_VENDOR = -mod vendor|g' workday-counter-$(VERSION)/Makefile
	#make -C workday-counter-$(VERSION) vendor
	tar --exclude .git -cJf workday-counter-$(VERSION).tar.xz workday-counter-$(VERSION)
	rm -rf workday-counter-$(VERSION)
