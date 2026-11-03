SHELL := bash

# Build the binary
build:
	hack/build.sh

# Build all release artifacts except android
release:
	hack/release.sh

# Run linter
lint:
	golangci-lint run -v --timeout 300s

# Run unit-tests
test:
	go test -v -coverprofile=coverprofile.out ./cmd/... ./pkg/...

# Generate coverage profile
coverprofile:
	hack/coverprofile.sh

# Format Go code
fmt:
	gofmt -s -w ./cmd ./pkg ./tests

# Validate that the generated files are up to date
validate:
	hack/validate.sh

# Validate the appstream metainfo file
validate-metainfo:
	appstreamcli validate io.github.heathcliff26.kde-trayrotate.metainfo.xml

# Update project dependencies
update-deps:
	hack/update-deps.sh

# Scan code for vulnerabilities using gosec
gosec:
	gosec ./...

# Build rpm with code in current workdir using packit
packit:
	packit build locally

# Build rpm of upstream code using packit + mock
packit-mock:
	packit build in-mock --resultdir tmp
	rm *.src.rpm

# Clean up build artifacts and temporary files
clean:
	hack/clean.sh

# Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t
	@echo ""
	@echo "Run 'make <target>' to execute a specific target."

.PHONY: \
	build \
	release \
	lint \
	test \
	coverprofile \
	fmt \
	validate \
	validate-metainfo \
	update-deps \
	gosec \
	packit \
	packit-mock \
	clean \
	help \
	$(NULL)
