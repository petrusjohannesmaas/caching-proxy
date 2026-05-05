# Caching Proxy

## 📡 Overview

A caching server built with Go. It forwards requests to an origin server, caches responses to disk, and improves performance on repeated requests.

* **Disk-based caching** — responses persist across restarts
* **Fast performance using Go's concurrency**
* **Simple flat-file storage** — each cached path becomes a file
* **Compiles to binary**

> 💡 **Prerequisites:** Make sure you have Go installed on your machine if you want to recompile the binary.

## 🧠 How It Works

| Command                                       | Description                        |
| --------------------------------------------- | ---------------------------------- |
| `./caching-proxy --port 3000 --origin http://dummyjson.com`   | Starts the Proxy      |
| `./caching-proxy --clear-cache`               | Manually clear the cache          |
| `./caching-proxy --origin http://example.com --cache-dir /var/cache/proxy` | Custom cache directory |

### Send a request to:

```
http://localhost:3000/products
```

You should receive the response from the origin and subsequent requests will hit the cache.

### Clone the Repository

```sh
git clone https://github.com/petrusjohannesmaas/roadmap.sh
cd caching-proxy
```

### Build the Project

```sh
go build -o caching-proxy proxy.go
```

This will compile a binary named `caching-proxy` in the project folder.

## 📦 Recommended Libraries

Consider these packages to enhance or optimize your proxy:

| Package                              | Description                         |
| ------------------------------------ | ----------------------------------- |
| [`fasthttp`](https://github.com/valyala/fasthttp) | High-performance HTTP server/client |
| [`fiber`](https://github.com/gofiber/fiber)   | Express.js-style web framework for Go |
| [`bbolt`](https://github.com/etcd-io/bbolt)   | Embedded key-value store            |

## 📈 Future Enhancements

* Add TTL-based cache expiration
* Add max cache size with LRU-style eviction
* Improve logging & error handling
* Support more HTTP methods (e.g., POST, PUT)
* Secure endpoint for manual cache invalidation

## 📄 License

MIT License © [Petrus Johannes Maas](https://github.com/petrusjohannesmaas)
