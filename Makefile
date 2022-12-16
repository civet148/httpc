#SHELL=/usr/bin/env bash

CLEAN:=
BINS:=
DATE_TIME=`date +'%Y%m%d %H:%M:%S'`
COMMIT_ID=`git rev-parse --short HEAD`
MANAGER_DIR=${PWD}
CONSOLE_CODE=/tmp/node-console

manager:
	rm -f httpc
	go mod tidy \
	&& go build -ldflags "-s -w -X 'main.BuildTime=${DATE_TIME}' -X 'main.GitCommit=${COMMIT_ID}'" -o httpc cmd/main.go
.PHONY: manager
BINS+=httpc

clean:
	rm -rf $(CLEAN) $(BINS)
.PHONY: clean
