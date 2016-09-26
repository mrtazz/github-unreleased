#
# some housekeeping tasks
#

export GO15VENDOREXPERIMENT = 1

# variable definitions
NAME := github-unreleased
DESC := experimental linter for Makefiles
PREFIX ?= usr/local
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDDATE := $(shell date -u +"%B %d, %Y")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
PKG_RELEASE ?= 1
PROJECT_URL := "https://github.com/mrtazz/$(NAME)"
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.buildTime=$(BUILDTIME)' \
           -X 'main.builder=$(BUILDER)' \
           -X 'main.goversion=$(GOVERSION)'

PACKAGES := $(shell find ./* -type d | grep -v vendor)


CMD_SOURCES := $(shell find cmd -name main.go)
TARGETS := $(patsubst cmd/%/main.go,%,$(CMD_SOURCES))
MAN_SOURCES := $(shell find man -name "*.md")
MAN_TARGETS := $(patsubst man/man1/%.md,%,$(MAN_SOURCES))

INSTALLED_TARGETS = $(addprefix $(PREFIX)/bin/, $(TARGETS))
INSTALLED_MAN_TARGETS = $(addprefix $(PREFIX)/share/man/man1/, $(MAN_TARGETS))

# source, dependency and build definitions
DEPDIR = .d
MAKEDEPEND = echo "$@: $$(go list -f '{{ join .Deps "\n" }}' $< | awk '/github/ { gsub(/^github.com\/[a-z_-]*\/[a-z_-]*\//, ""); printf $$0"/*.go " }')" > $(DEPDIR)/$@.d

$(DEPDIR)/%.d: ;
.PRECIOUS: $(DEPDIR)/%.d

$(DEPDIR):
	install -d $@

-include $(patsubst %,$(DEPDIR)/%.d,$(TARGETS))

%: cmd/%/main.go $(DEPDIR) $(DEPDIR)/%.d
	$(MAKEDEPEND)
	go build -ldflags "$(LDFLAGS)" -o $@ $<

%.1: man/man1/%.1.md
	sed "s/REPLACE_DATE/$(BUILDDATE)/" $< | pandoc -s -t man -o $@

all: $(TARGETS) $(MAN_TARGETS)
.DEFAULT_GOAL:=all

# development tasks
test:
	go test -v $$(go list ./... | grep -v /vendor/)

coverage:
	@echo "mode: set" > cover.out
	@for package in $(PACKAGES); do \
		if ls $${package}/*.go &> /dev/null; then  \
		go test -coverprofile=$${package}/profile.out $${package}; fi; \
		if test -f $${package}/profile.out; then \
	 	cat $${package}/profile.out | grep -v "mode: set" >> cover.out; fi; \
	done
	@-go tool cover -html=cover.out -o cover.html

benchmark:
	@echo "Running tests..."
	@go test -bench=. ${NAME}

govendor:
	    go get -u github.com/kardianos/govendor

# install tasks
$(PREFIX)/bin/%: %
	install -d $$(dirname $@)
	install -m 755 $< $@

$(PREFIX)/share/man/man1/%: %
	install -d $$(dirname $@)
	install -m 644 $< $@

install: $(INSTALLED_TARGETS) $(INSTALLED_MAN_TARGETS)

local-install:
	$(MAKE) install PREFIX=usr/local

# packaging tasks
packages: local-install rpm deb

deploy-packages: packages
	package_cloud push mrtazz/$(NAME)/el/7 *.rpm
	package_cloud push mrtazz/$(NAME)/debian/wheezy *.deb
	package_cloud push mrtazz/$(NAME)/ubuntu/trusty *.deb


rpm: $(SOURCES)
	  fpm -t rpm -s dir \
    --name $(NAME) \
    --version $(VERSION) \
		--description "$(DESC)" \
    --iteration $(PKG_RELEASE) \
    --epoch 1 \
    --license MIT \
    --maintainer "Daniel Schauenberg <d@unwiredcouch.com>" \
    --url $(PROJECT_URL) \
    --vendor mrtazz \
    usr

deb: $(SOURCES)
	  fpm -t deb -s dir \
    --name $(NAME) \
    --version $(VERSION) \
		--description "$(DESC)" \
    --iteration $(PKG_RELEASE) \
    --epoch 1 \
    --license MIT \
    --maintainer "Daniel Schauenberg <d@unwiredcouch.com>" \
    --url $(PROJECT_URL) \
    --vendor mrtazz \
    usr


# clean up tasks
clean-deps:
	$(RM) -r $(DEPDIR)

clean: clean-docs clean-deps
	$(RM) -r ./usr
	$(RM) $(TARGETS)

clean-docs:
	$(RM) $(MAN_TARGETS)

pizza:
	@echo ""
	@echo "🍕 🍕 🍕 🍕 🍕 🍕   make.pizza 🍕 🍕 🍕 🍕 🍕 🍕 "
	@echo ""
	@echo "https://twitter.com/mrb_bk/status/760636493710983168"
	@echo ""

.PHONY: all test rpm deb install local-install packages govendor coverage clean-deps clean clean-docs pizza
