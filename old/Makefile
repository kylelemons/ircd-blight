VERSION=$( git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD )

EMBED_VERSION=-ldflags "-X github.com/kylelemons/ircd-blight/old/ircd/core/REPO_VERSION=$(VERSION)"

.PHONY : generate
generate :
	go generate ./...

.PHONY : build
build : generate
	go build $(EMBED_VERSION) ./...

.PHONY : install
install : generate
	go install $(EMBED_VERSION) ./cmd/*
