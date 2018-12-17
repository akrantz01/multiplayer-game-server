# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get

# Naming of the output
BINARY_NAME=server
LINUX_SUFFIX=linux
MACOS_SUFFIX=macos
WIN_SUFFIX=windows
ARCH64_SUFFIX=amd64
ARCH32_SUFFIX=i386

# Where to put stuff
OUT_FOLDER=./dist

# Help info (default target)
help:
	@echo "MMOS Makefile Help:"
	@echo ""
	@echo "Standard build targets:"
	@echo "	all			builds binary for every OS and architecture"
	@echo "	build			builds binary for your OS and architecture (via GOOS and GOARCH)"
	@echo "	clean			clean the build environment"
	@echo "	deps			get the required dependencies"
	@echo ""
	@echo "OS/Architecture build targets:"
	@echo "	build-linux-amd64	builds binary for linux on amd64 architecture"
	@echo "	build-linux-i386	builds binary for linux on i386 architecture"
	@echo "	build-windows-amd64	builds binary for windows on amd64 architecture"
	@echo "	build-windows-i386	builds binary for windows on i386 architecture"
	@echo "	build-macos		builds binary for macos on amd64 architecture"

# General targets
all: clean deps build-linux-amd64 build-linux-i386 build-windows-amd64 build-windows-i386 build-macos

build: clean deps
	$(GOBUILD) -o $(OUT_FOLDER)/$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -rf $(OUT_FOLDER)

deps:
	$(GOGET) github.com/gorilla/websocket
	$(GOGET) github.com/gorilla/handlers
	$(GOGET) gopkg.in/yaml.v2

# OS/Arch specific build commands
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUT_FOLDER)/$(BINARY_NAME)-$(LINUX_SUFFIX)-$(ARCH64_SUFFIX)

build-linux-i386:
	GOOS=linux GOARCH=386 $(GOBUILD) -o $(OUT_FOLDER)/$(BINARY_NAME)-$(LINUX_SUFFIX)-$(ARCH32_SUFFIX)

build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(OUT_FOLDER)/$(BINARY_NAME)-$(WIN_SUFFIX)-$(ARCH64_SUFFIX)

build-windows-i386:
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(OUT_FOLDER)/$(BINARY_NAME)-$(WIN_SUFFIX)-$(ARCH32_SUFFIX)

build-macos:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(OUT_FOLDER)/$(BINARY_NAME)-$(MACOS_SUFFIX)