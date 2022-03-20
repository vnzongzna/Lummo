package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	kv "github.com/vnzongzna/lummo/pkg/kv/v1"
	"github.com/vnzongzna/lummo/pkg/metrics"
)

func main() {
	kv := kv.Init()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(metrics.PrometheusMiddleware)

	r.Get("/get/{key}", kv.Get)
	r.Post("/set", kv.Set)
	r.Get("/search", kv.Search)
	r.Handle("/prometheus", metrics.Handler())

	log.Fatal(http.ListenAndServe(":8088", r))
}
