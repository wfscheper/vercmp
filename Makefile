COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/wfscheper/vercmp/vercmp.CommitHash=${COMMIT_HASH} -X github.com/wfscheper/vercmp/vercmp.BuildDate=${BUILD_DATE}"
PACKAGES = $(shell govendor list -no-status +local | sed 's/github.com.wfscheper.vercmp/./')

all: gitinfo

install: install-gitinfo

help:
	echo ${COMMIT_HASH}
	echo ${BUILD_DATE}

gitinfo:
	go build ${LDFLAGS} vercmp.go

install-gitinfo:
	go install ${LDFLAGS} ./...

no-git-info:
	go build vercmp.go

govendor:
	go get -u github.com/kardianos/govendor
	go install github.com/kardianos/govendor
	govendor sync github.com/wfscheper/vercmp

check: fmt vet test test-race

cyclo:
	@for d in `govendor list -no-status +local | sed 's/github.com.wfscheper.vercmp/./'` ; do \
		if [ "`gocyclo -over 20 $$d | tee /dev/stderr`" ]; then \
			echo "^ cyclomatic complexity exceeds 20, refactor the code!" && echo && exit 1; \
		fi \
	done

fmt:
	@for d in `govendor list -no-status +local | sed 's/github.com.wfscheper.vercmp/./'` ; do \
		if [ "`gofmt -l $$d/*.go | tee /dev/stderr`" ]; then \
			echo "^ improperly formatted go files" && echo && exit 1; \
		fi \
	done

lint:
	@for d in `govendor list -no-status +local | sed 's/github.com.wfscheper.vercmp/./'` ; do \
		if [ "`golint $$d | tee /dev/stderr`" ]; then \
			echo "^ golint errors!" && echo && exit 1; \
		fi \
	done

get:
	go get -v -t ./...

test:
	govendor test +local

test-race:
	govendor test -race +local

vet:
	@if [ "`govendor vet +local | tee /dev/stderr`" ]; then \
		echo "^ go vet errors!" && echo && exit 1; \
	fi

test-cover-html:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		govendor test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out
