mitm
====

[![Build Status](https://travis-ci.org/fortytw2/mitm.svg?branch=master)](https://travis-ci.org/fortytw2/mitm) [![GoDoc](https://godoc.org/github.com/fortytw2/mitm?status.svg)](https://godoc.org/github.com/fortytw2/mitm)

TLS MITM Detection based on the work done in `caddy`, extracted as an 
easily importable and seperately testable go library for usage in 
independent go programs. 

This capability is based on research done by Durumeric,
Halderman, et. al. in "The Security Impact of HTTPS Interception" (NDSS '17) - https://jhalderm.com/pub/papers/interception-ndss17.pdf

Full credit to [caddy](https://github.com/mholt/caddy) and all contributors for the implementation here, I've just tidied it up for use in other projects
and given it a clean public API free of `caddy/` imports and external 
dependencies.

Usage
-----

Pass your `net.Listener` and `*tls.Config` to `mitm`, then check to see if your connections are MITM-ed with `mitm.Check(r)`.

```go
package main

import (
	"log"
	"net"
	"net/http"

	"github.com/fortytw2/mitm"
)

func main() {
	l, err := net.Listen("localhost", "8080")
	if err != nil {
		panic(err)
	}

	h := mitm.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mitm.Check(r) {
			panic("oh no, you've been MITM-ed")
		}
		w.Write([]byte("all good!"))
	}), l, nil, false)

	log.Fatal(http.Serve(l, h))
}
```

Installation
------------

```bash
go get -u github.com/fortytw2/mitm/...
```

License
-------

Apache, see LICENSE
