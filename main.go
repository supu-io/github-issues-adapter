package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/supu-io/payload"
)

func setupRouter() *martini.ClassicMartini {
	m := martini.Classic()
	m.Post("/move", Move)

	return m
}

func main() {
	m := setupRouter()
	m.Run()
}

// Move moves that issue on the
func Move(w http.ResponseWriter, r *http.Request) string {
	decoder := json.NewDecoder(r.Body)
	var p payload.Payload
	decoder.Decode(&p)
	doMove(p, nil)

	return "success"
}
