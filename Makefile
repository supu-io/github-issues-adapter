deps:
	go get -u github.com/google/go-github/github
	go get -u golang.org/x/oauth2
	go get -u github.com/supu-io/payload
dev-deps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/go-martini/martini
	go get -u github.com/smartystreets/goconvey/convey
build:
	go build
test:
	go test
lint:
	golint
