deps:
	go get -u github.com/google/go-github/github
	go get -u golang.org/x/oauth2
	go get -u github.com/supu-io/payload
build:
	go build
test:
	go test
