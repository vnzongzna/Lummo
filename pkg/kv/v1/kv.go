package kv

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

type KVStore struct {
	data map[string]string
	// anti pattern in golang, but why not
	sync.Mutex
}

// initialize data-structure so we can use it
func Init() *KVStore {
	return &KVStore{
		data:  make(map[string]string),
		Mutex: sync.Mutex{},
	}
}

func (kv *KVStore) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	// set default headers
	w.Header().Set("Content-Type", "application/json")

	// acquire lock before reading from common map
	kv.Lock()
	defer kv.Unlock()

	// if the key exist, convert it into a json
	if value, ok := kv.data[key]; ok {
		marshalled, err := json.Marshal(&struct {
			Value string `json:"value"`
		}{
			Value: value,
		})

		// in case marshalling failed for any reason, give proper response code
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("failed to unmarshal result")
		}

		w.Write(marshalled)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (kv *KVStore) Set(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	t := make(map[string]string)

	if err := decoder.Decode(&t); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Failed to decode ", t)
	} else {
		// only write data if decoder had no error
		// acquire lock to write data to common map
		kv.Lock()
		defer kv.Unlock()
		for k, v := range t {
			kv.data[k] = v
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

type searchResult struct {
	Keys []string `json:"keys"`
}

func (kv *KVStore) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	result := make([]string, 0)

	w.Header().Set("Content-Type", "application/json")

	kv.Lock()
	defer kv.Unlock()

	if query.Has("prefix") {
		prefix := query.Get("prefix")
		for k := range kv.data {
			if strings.HasPrefix(k, prefix) {
				result = append(result, k)
			}
		}
		w.WriteHeader(http.StatusOK)
	} else if query.Has("suffix") {
		suffix := query.Get("suffix")
		for k := range kv.data {
			if strings.HasSuffix(k, suffix) {
				result = append(result, k)
			}
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	result_marshalled, err := json.Marshal(&searchResult{Keys: result})
	if err != nil {
		log.Println("marshalling failed for", result)
	}
	w.Write(result_marshalled)
}
