TEST?=./...
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOLINT_TARGETS   ?= $$(golint github.com/yamamoto-febc/terraform-provisioner-vnc/vnc | grep -v 'underscores in Go names' | tee /dev/stderr)
CURRENT_VERSION = $(shell gobump show -r vnc/)
BUILD_LDFLAGS = "-s -w \
	  -X github.com/yamamoto-febc/terraform-provisioner-vnc/vnc.Revision=`git rev-parse --short HEAD`"
export GO111MODULE=on

default: test build

.PHONY: tools
tools:
	GO111MODULE=off go get -u github.com/motemen/gobump/cmd/gobump
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get -u golang.org/x/lint/golint

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean
	OS="`go env GOOS`" ARCH="`go env GOARCH`" ARCHIVE= BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

build-x: build-darwin build-windows build-linux build-bsd shasum

build-darwin: bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_darwin-386.zip bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_darwin-amd64.zip

build-windows: bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_windows-386.zip bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_windows-amd64.zip

build-linux: bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_linux-386.zip bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_linux-amd64.zip bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_linux-arm.zip

build-bsd: bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_openbsd-386.zip bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_openbsd-amd64.zip bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_openbsd-arm.zip

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_darwin-386.zip:
	OS="darwin"  ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_darwin-amd64.zip:
	OS="darwin"  ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_windows-386.zip:
	OS="windows" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_windows-amd64.zip:
	OS="windows" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_linux-386.zip:
	OS="linux"   ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_linux-amd64.zip:
	OS="linux"   ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_linux-arm.zip:
	OS="linux"   ARCH="arm" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_openbsd-386.zip:
	OS="openbsd" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_openbsd-amd64.zip:
	OS="openbsd" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provisioner-vnc_$(CURRENT_VERSION)_openbsd-arm.zip:
	OS="openbsd" ARCH="arm" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

shasum:
	(cd bin/; shasum -a 256 * > terraform-provisioner-vnc_$(CURRENT_VERSION)_SHA256SUMS)

test: fmt
	TF_ACC=  go test $(TEST) -v $(TESTARGS) -timeout=30s -parallel=4 ; \

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 240m ; \

golint: goimports
	test -z "$(GOLINT_TARGETS)"

goimports: fmt
	goimports -w $(GOFMT_FILES)

fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: default test testacc fmt fmtcheck

.PHONY: bump-patch bump-minor bump-major version
bump-patch:
	gobump patch -w vnc

bump-minor:
	gobump minor -w vnc

bump-major:
	gobump major -w vnc

version:
	gobump show -r vnc

git-tag:
	git tag `gobump show -r vnc`
