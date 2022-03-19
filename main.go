package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	kv "github.com/vnzongzna/lummo/kv/v1"
)

func main() {
	kv := kv.Init()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/get/{key}", kv.Get)
	r.Post("/set", kv.Set)
	r.Get("/search", kv.Search)

	log.Fatal(http.ListenAndServe(":80", r))
}
