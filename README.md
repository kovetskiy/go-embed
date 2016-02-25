## Go Embed
Generates go code to embed resource files into your library or executable.
This is more suitable for web servers as it gzip compresses all the files
automatically and computes the hash so that it can be used for caching the
assets in the frontends.

```bash
go-embed v1.0.0
Generates go code to embed resource files into your library or executable

  Usage:
    -input  string  The path to the folder containing the assets
    -output string  The output filename

  example:
    go-embed -input public/ -output assets/main.go
```

You can use this to embed your css, js and images into a single executable.

This is similar to [go-bindata](https://github.com/jteeuwen/go-bindata).

This is similar to [pony-embed](https://github.com/pyros2097/pony-embed).

This is similar to [rust-embed](https://github.com/pyros2097/rust-embed).

## Installation
```
go get github.com/pyros2097/go-embed
```

## Documentation
You can directly access your files as constants from the assets package or
you can use this func to serve all files stored in your assets folder which is useful for webservers and has gzip compression and hash creation for caching. Just see the example as to how same caching and compression works in
production and development.
```go
assets.IndexHTML // direct access
assets.Asset(base, path) (data, hash, contentType, error)
```

## Examples
A simple http server which serves its resources directly.

Navigate to the example folder and run these commands,

To see it in action in development,
`go run main.go`

To see it in action in production,
```
cd examples
go-embed -input public/ -output assets/main.go
go run main.go
git stash
git stash drop
```

```go
package main

import (
  "net/http"

  "github.com/pyros2097/go-embed/examples/assets"
)

func main() {
  http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    println("GET " + req.URL.Path)
    data, hash, contentType, err := assets.Asset("public/", req.URL.Path)
    if err != nil {
      data, hash, contentType, err = assets.Asset("public", "/index.html")
      if err != nil {
        data = []byte(err.Error())
      }
    }
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
```

For development mode you can use this as a template file for running your
server and loading the assets directly from the filesystem and in release
this file would be rewritten by the go-embed tool to contain the actual file
data and after you run your build command you might need to checkout this file
again. This file can be safely committed to your git repository.
```go
package assets

import (
  "bytes"
  "compress/gzip"
  "crypto/md5"
  "encoding/hex"
  "io/ioutil"
  "mime"
  "path/filepath"
)

func init() {
  mime.AddExtensionType(".ico", "image/x-icon")
  mime.AddExtensionType(".eot", "font/eot")
  mime.AddExtensionType(".tff", "font/tff")
  mime.AddExtensionType(".woff", "application/font-woff")
  mime.AddExtensionType(".woff2", "application/font-woff")
}

// Asset Gets the file from system if debug otherwise gets it from the stored
// data returns the data, the md5 hash of its content and its content type and
// and error if it is not found
// In production the base parameter is ignored
// use it for developement as the relativepath/to/your/public/folder
// ex: ui/myproject/src/public
func Asset(base, path string) ([]byte, string, string, error) {
  var b bytes.Buffer
  file := base + path
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return nil, "", "", err
  }
  if data != nil {
    w := gzip.NewWriter(&b)
    w.Write(data)
    w.Close()
    data = b.Bytes()
  }
  sum := md5.Sum(data)
  return data, hex.EncodeToString(sum[1:]), mime.TypeByExtension(filepath.Ext(file)), nil
}
```
This is how I use it in my make file,
```bash
cd ui/$(Project) && npm install && npm run build
go-embed -input ui/$(Project)/src/public -output assets/main.go
go build -o $(Project) -tags '$(Project)' $(Project).go
git checkout assets/main.go
```

Go Gophers!

The power is yours!
