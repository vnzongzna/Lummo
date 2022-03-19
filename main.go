package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	kv "github.com/vnzongzna/lummo/kv/v1"
)

func main() {
	kv := kv.Init()
	r := mux.NewRouter()

	r.HandleFunc("/get/{key}", kv.Get).Methods("GET").Headers("Content-Type", "application/json")
	r.HandleFunc("/set", kv.Set).Methods("POST").Headers("Content-Type", "application/json")
	r.HandleFunc("/search", kv.Search).Methods("GET").Headers("Content-Type", "application/json")

	log.Fatal(http.ListenAndServe(":80", r))
}
