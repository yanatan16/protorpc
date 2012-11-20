# Makefile for a go project
#
# Author: Jon Eisen
# 	site: joneisen.me
# 	
# Targets:
# 	all: Builds the code
# 	build: Builds the code
# 	fmt: Formats the source files
# 	clean: cleans the code
# 	install: Installs the code to the GOPATH
# 	iref: Installs referenced projects
#	test: Runs the tests
#	
#  Blog post on it: http://joneisen.me/post/25503842796
#

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

# Package lists
TOPLEVEL_PKG := protorpc

# List building
CMD_LIST = broker protoc-gen-go-rpc
ALL_LIST = $(foreach int, $(CMD_LIST), cmd/$(int))

BUILD_LIST = $(foreach int, $(ALL_LIST), $(int)_build)
CLEAN_LIST = $(foreach int, $(ALL_LIST), $(int)_clean)
INSTALL_LIST = $(foreach int, $(ALL_LIST), $(int)_install)
IREF_LIST = $(foreach int, $(ALL_LIST), $(int)_iref)
TEST_LIST = $(foreach int, $(ALL_LIST), $(int)_test)
FMT_TEST = $(foreach int, $(ALL_LIST), $(int)_fmt)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)

all: build
build: $(BUILD_LIST)
	$(GOBUILD) $(TOPLEVEL_PKG)
clean: $(CLEAN_LIST)
	$(GOCLEAN) $(TOPLEVEL_PKG)
install: $(INSTALL_LIST)
	$(GOINSTALL) $(TOPLEVEL_PKG)
test: $(TEST_LIST)
	$(GOTEST) $(TOPLEVEL_PKG)
iref: $(IREF_LIST)
	$(GODEP) $(TOPLEVEL_PKG)
fmt: $(FMT_TEST)
	$(GOFMT) .

$(BUILD_LIST): %_build: %_fmt %_iref
	$(GOBUILD) $(TOPLEVEL_PKG)/$*
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(TOPLEVEL_PKG)/$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*
$(IREF_LIST): %_iref:
	$(GODEP) $(TOPLEVEL_PKG)/$*
$(TEST_LIST): %_test:
	$(GOTEST) $(TOPLEVEL_PKG)/$*
$(FMT_TEST): %_fmt:
	$(GOFMT) ./$*

init:
	ln -s `pwd` $(GOPATH)/src/protorpc
	go install protorpc