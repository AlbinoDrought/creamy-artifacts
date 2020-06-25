.PHONY: default
default: deps

.PHONY: deps
deps:
	go get -d -v

.PHONY: build
build: deps
	go build

.PHONY: test
test: deps
	go test

.PHONY: run
run: build
	./creamy-artifacts

.PHONY: install
install: deps
	go install

.PHONY: uninstall
uninstall:
	go clean -x -i

.PHONY: clean
clean:
	go clean -x

.PHONY: image
image:
	docker build -t albinodrought/creamy-artifacts .
