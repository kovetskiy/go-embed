package main

import (
	"net/http"

	"github.com/pyros2097/go-embed/examples/assets"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		println("GET " + req.URL.Path)
		data, hash, contentType := assets.Asset(req.URL.Path)
		res.Header().Set("Content-Encoding", "gzip")
		res.Header().Set("Content-Type", contentType)
		res.Header().Add("Cache-Control", "public, max-age=31536000")
		res.Header().Add("ETag", hash)
		if req.Header.Get("If-None-Match") == hash {
			res.WriteHeader(http.StatusNotModified)
		} else {
			res.WriteHeader(http.StatusOK)
			_, err := res.Write(data)
			if err != nil {
				panic(err)
			}
		}
	})
	println("Server running on 127.0.0.1:3000")
	http.ListenAndServe(":3000", nil)
}
