.PHONY: tests

tests:
	find . -maxdepth 4 -iname '*_test.go' | xargs -n 1 dirname | xargs godep go test

bin: subcommands Makefile
	find subcommands -maxdepth 3 -iname '*.go' \
		| xargs -n 1 dirname \
		| sort | uniq | \
		xargs -I % go install github.com/tbd-ci/tbd/%
