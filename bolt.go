package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	readPath = regexp.MustCompile("^/u/(\\d+)$")
)

type bolt struct {
	data []string
}

func (h *bolt) persist() {
	if dataDir == "" {
		return
	}

	byts, err := json.Marshal(h.data)
	if err != nil {
		log.Fatalf("Quitting, cannot marshal data err=%v", err)
	}

	err = ioutil.WriteFile(dataFilePath, byts, 0644)
	if err != nil {
		log.Fatalf("Quitting, cannot write data err=%v", err)
	}
}

func (h *bolt) store(w http.ResponseWriter, r *http.Request) {
	urls, ok := r.URL.Query()["url"]
	if !ok {
		http.Error(w, "Missing parameter: url", 400)
		return
	}

	if len(urls) != 1 {
		http.Error(w, "Expected single parameter: url", 400)
		return
	}

	if urls[0] == "" {
		http.Error(w, "Missing parameter: url", 400)
		return
	}

	h.data = append(h.data, urls[0])
	h.persist()

	w.Write([]byte(fmt.Sprintf("/u/%x", len(h.data))))
}

func (h *bolt) read(w http.ResponseWriter, r *http.Request) {
	matches := readPath.FindAllStringSubmatch(r.URL.Path, -1)
	if len(matches) != 1 || len(matches[0]) != 2 {
		log.Printf("Request for invalid path [%#v] returning 404.\n", r)
		http.Error(w, "Not found", 404)
		return
	}

	idx, err := strconv.Atoi(matches[0][1])
	if err != nil || idx > len(h.data) {
		log.Printf("Failed to convert index %#v, err=%v.\n", matches[0][1], err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, h.data[idx-1], http.StatusFound)
}

func (h *bolt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/s/"):
		h.store(w, r)
	case strings.HasPrefix(r.URL.Path, "/u/"):
		h.read(w, r)
	default:
		log.Printf("Ignoring request as 404: %#v\n", r)
		http.Error(w, "Not found.", 404)
	}
}

func serve(b *bolt, addr string) {
	mux := http.NewServeMux()
	mux.Handle("/", b)
	log.Printf("bolt started at [%v] with %v entries.", addr, len(b.data))
	log.Fatal(http.ListenAndServe(addr, mux))
}
