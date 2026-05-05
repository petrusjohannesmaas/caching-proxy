package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	port          int
	origin        string
	clearCache    bool
	cacheDir      string
)

func fetchAndCache(w http.ResponseWriter, req *http.Request) {
	url := origin + req.URL.Path
	log.Printf("Fetching from origin: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	cachePath := filepath.Join(cacheDir, req.URL.Path)
	os.MkdirAll(filepath.Dir(cachePath), 0755)
	os.WriteFile(cachePath, data, 0644)
	w.Header().Set("X-Cache", "MISS")
	w.Write(data)
}

func handler(w http.ResponseWriter, req *http.Request) {
	cachePath := filepath.Join(cacheDir, req.URL.Path)
	if data, err := os.ReadFile(cachePath); err == nil {
		w.Header().Set("X-Cache", "HIT")
		w.Write(data)
		log.Printf("Served from cache: %s", req.URL.Path)
	} else {
		fetchAndCache(w, req)
	}
}

func main() {
	flag.IntVar(&port, "port", 3000, "Port for proxy server")
	flag.StringVar(&origin, "origin", "", "Origin server URL")
	flag.StringVar(&cacheDir, "cache-dir", "/tmp/caching-proxy", "Cache directory")
	flag.BoolVar(&clearCache, "clear-cache", false, "Clears the cache and exits")

	flag.Parse()

	if clearCache {
		os.RemoveAll(cacheDir)
		fmt.Println("Cache cleared!")
		return
	}

	if origin == "" {
		fmt.Println("Usage: caching-proxy --port <number> --origin <url>")
		return
	}

	os.MkdirAll(cacheDir, 0755)

	http.HandleFunc("/", handler)
	fmt.Printf("Caching proxy server running on port %d, forwarding to %s\n", port, origin)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
