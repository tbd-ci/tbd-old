
tests:
	find . -maxdepth 3 -iname *_test.go | xargs -n 1 dirname | xargs godep go test
