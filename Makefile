
.PHONY: build
build:
	go install .

.PHONY: build-lint
build-lint:
	go get -u github.com/alecthomas/gometalinter && gometalinter --install

.PHONY: lint
lint: build-lint
	gometalinter ./... --vendor --config gometalinter.conf --exclude='.*(easyjson.go|pb.go|test.go)'

.PHONY: test
test:
	go test -v $(shell go list ./...)

