package kv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type KVStore struct {
	data map[string]string
	sync.Mutex
}

func Init() KVStore {
	return KVStore{
		data:  make(map[string]string),
		Mutex: sync.Mutex{},
	}
}

func (kv *KVStore) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("get", r.RequestURI)
	vars := mux.Vars(r)
	kv.Lock()
	defer kv.Unlock()
	if value, ok := kv.data[vars["key"]]; ok {
		w.WriteHeader(http.StatusOK)
		marshalled, err := json.Marshal(struct {
			Value string `json:"value"`
		}{
			Value: value,
		})
		if err != nil {
			log.Println("failed to unmarshal result")
		}
		fmt.Fprintln(w, string(marshalled))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (kv *KVStore) Set(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	t := make(map[string]string)

	if err := decoder.Decode(&t); err != nil {
		log.Println("Failed to decode ", t)
	}
	log.Println("set", t)
	kv.Lock()
	defer kv.Unlock()
	for k, v := range t {
		kv.data[k] = v
	}
}

func (kv *KVStore) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	result := make([]string, 0)

	log.Println("search", query)
	kv.Lock()
	defer kv.Unlock()

	if query.Has("prefix") {
		w.WriteHeader(http.StatusOK)
		prefix := query.Get("prefix")
		for k := range kv.data {
			if strings.HasPrefix(k, prefix) {
				result = append(result, k)
			}
		}
	} else if query.Has("suffix") {
		w.WriteHeader(http.StatusOK)
		suffix := query.Get("prefix")
		for k := range kv.data {
			if strings.HasSuffix(k, suffix) {
				result = append(result, k)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result_marshalled, err := json.Marshal(struct {
		Keys []string `json:"keys"`
	}{
		Keys: result,
	})
	if err != nil {
		log.Println("marshalling failed for", result)
	}
	fmt.Fprintf(w, string(result_marshalled))
}
