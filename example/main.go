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
