package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-github/github"
	"github.com/supu-io/payload"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *github.Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	client.UploadURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestIssuesService_Update(t *testing.T) {
	setup()
	defer teardown()

	input := []string{"doing"}

	mux.HandleFunc("/repos/o/r/issues/1/labels/todo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels/doing", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels/review", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels/uat", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels/done", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels", func(w http.ResponseWriter, r *http.Request) {
		v := new([]string)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(*v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `[{"name":"doing"}]`)
	})

	id := "o/r/1"
	issue := payload.Issue{ID: &id}
	from := "todo"
	to := "doing"
	trans := payload.Transition{From: &from, To: &to}
	token := "token"
	github := payload.Github{Token: &token}
	config := payload.Config{Github: &github}
	status := []string{"todo", "doing", "review", "uat", "done"}

	p := payload.Payload{
		Issue:      &issue,
		Transition: &trans,
		Config:     &config,
		Status:     &status,
	}
	err := doMove(p, client)

	if err != nil {
		t.Errorf(err.Error())
	}

}
