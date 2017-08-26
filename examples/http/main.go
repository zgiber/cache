package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/zgiber/cache"
	"goji.io"
	"goji.io/pat"
)

var (
	httpPort               = "8000"
	memory                 = cache.NewMemCache()
	defaultCacheExpiration = 10 * time.Minute
)

// retrieve a resource with a key
func handleGetResource(w http.ResponseWriter, r *http.Request) {
	key := pat.Param(r, "id")
	value, err := memory.Fetch(key)
	if err == nil {
		w.Write(value)
	}

	// here you'd want to retrieve
	// the value from some persistent
	// storage, and update the cache
	// if you want to refresh the expiry.
}

// upsert a resource with a key
func handlePutResource(w http.ResponseWriter, r *http.Request) {
	key := pat.Param(r, "id")
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// here you'd want to upsert
	// the value to some persistent
	// storage and return error if
	// that fails.

	err = memory.Set(key, value, defaultCacheExpiration)
	_ = err // handle this error according to your taste
}

func main() {
	m := goji.NewMux()
	m.HandleFunc(pat.Get("/:id"), handleGetResource)
	m.HandleFunc(pat.Put("/:id"), handlePutResource)
	log.Println(http.ListenAndServe(":"+httpPort, m))
}
