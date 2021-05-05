# Cache

go-cache is a Go package which abstracts cache systems

Inspired from https://github.com/gookit/cache, but designed to be wrapped, similar to https://github.com/Shopify/go-storage.

# Requirements

- [Go 1.13+](http://golang.org/dl/)

# Installation

```console
$ go get github.com/Shopify/go-cache
```

# Usage

All caches in this package follow a simple [interface](client.go).

If you want, you may wrap the Client to inject functionality, like a circuit breaker, logging, or metrics.
