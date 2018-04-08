PREFIX			?= $(shell pwd)
PKGS 			:= $(shell go list ./... | grep -v /vendor/)
SHASUMCMD 		:= $(shell command -v sha1sum || command -v shasum; 2> /dev/null)
TARCMD 			:= $(shell command -v tar || command -v tar; 2> /dev/null)
GIT_REF 		:= $(shell git log --pretty=format:'%h' -n 1)
CURRENT_USER 	:= $(shell whoami)
VERSION 		:= $(shell cat ./.version)

export GIT_REF
export VERSION

all: format vet profanity test

ci: vet profanity test

new-install:
	@go get -u github.com/lib/lq

format:
	@echo "$(VERSION)/$(GIT_REF) >> formatting code"
	@go fmt $(PKGS)

vet:
	@echo "$(VERSION)/$(GIT_REF) >> vetting code"
	@go vet $(PKGS)

profanity:
	@echo "$(VERSION)/$(GIT_REF) >> profanity"
	@go run _bin/profanity/main.go

test:
	@echo "$(VERSION)/$(GIT_REF) >> running all tests"
	@go test $(PKGS)

increment-patch:
	@echo "Current Version $(VERSION)"
	@go run _bin/increment_version/main.go patch ./.version > ./NEW_VERSION
	@mv ./NEW_VERSION ./.version
	@cat ./.version

increment-minor:
	@echo "Current Version $(VERSION)"
	@go run _bin/increment_version/main.go minor ./.version > ./NEW_VERSION
	@mv ./NEW_VERSION ./.version
	@cat ./.version

increment-major:
	@echo "Current Version $(VERSION)"
	@go run _bin/increment_version/main.go major ./.version > ./NEW_VERSION
	@mv ./NEW_VERSION ./.version
	@cat ./.version

tag-version:
	git tag -f $(VERSION)

push-tags:
	git push -f --tags
