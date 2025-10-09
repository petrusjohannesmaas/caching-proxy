package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/golang-lru"
	"io"
	"log"
	"net/http"
)

var (
	port       int
	origin     string
	clearCache bool
	cache      *lru.Cache
)

// Fetch and cache responses
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

	cache.Add(req.URL.Path, data) // Add to cache
	w.Header().Set("X-Cache", "MISS")
	w.Write(data)
}

// Handle requests and serve cached responses
func handler(w http.ResponseWriter, req *http.Request) {
	if data, ok := cache.Get(req.URL.Path); ok {
		w.Header().Set("X-Cache", "HIT")
		w.Write(data.([]byte))
		log.Printf("Served from cache: %s", req.URL.Path)
	} else {
		fetchAndCache(w, req)
	}
}

func main() {
	// Define flags
	flag.IntVar(&port, "port", 3000, "Port for proxy server")
	flag.StringVar(&origin, "origin", "", "Origin server URL")
	flag.BoolVar(&clearCache, "clear-cache", false, "Clears the cache and exits")

	// Parse flags **only once**
	flag.Parse()

	// Handle cache clearing request separately
	if clearCache {
		cache, _ = lru.New(100) // Initialize cache
		cache.Purge()           // Clear all entries
		fmt.Println("Cache cleared!")
		return // Exit program immediately
	}

	// Ensure origin is provided
	if origin == "" {
		fmt.Println("Usage: caching-proxy --port <number> --origin <url>")
		return
	}

	// Initialize cache
	cache, _ = lru.New(100)

	// Set up HTTP server
	http.HandleFunc("/", handler)
	fmt.Printf("Caching proxy server running on port %d, forwarding to %s\n", port, origin)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil)) // Log errors if the server fails
}
