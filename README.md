## Go Embed
Generates go code to embed resource files into your library or executable.
This is more suitable for web servers as it gzip compresses all the files
automatically and computes the hash so that it can be used for caching the
assets in the frontends.

```bash
go-embed v0.1.0
Generates go code to embed resource files into your library or executable

  Usage:
    -input  string  The path to the folder containing the assets
    -output string  The output filename
    -tag    string  The tag to use for the generated package

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
## Requirements
* An `index.html` file in your input folder

## Documentation
You can directly access your files as constants from the assets package or
you can use this func to serve all files stored in your assets folder which is useful for webservers and has gzip compression and caching inbuilt. Just see the example as to how same caching and compression works in
production and development.
```go
assets.IndexHTML // direct access
assets.Asset(base, path) (data, hash, contentType)
```
The Asset func does not return an error like the rest of the resource embedding tools and that's because in normal SPA applications if a route is not found then you automatically redirect it to the root path.
Here in go-embed if a file or path is not found then we directly send the 
data for "index.html" which MUST be present in your input folder.
And the root path will always return the data for "index.html"

## Examples
A simple http server which serves its resources directly.

Navigate to the example folder and run these commands,

To see it in action in development,
`go run --tags dev main.go`

To see it in action in production,
`go run --tags prod main.go`

```go
package main

import (
  "net/http"

  "github.com/pyros2097/go-embed/examples/assets"
)

func main() {
  http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    println("GET " + req.URL.Path)
    data, hash, contentType := assets.Asset("", req.URL.Path)
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
data and after you run your build command you might need  to checkout this file
again. This file can be safely committed to your git repository.
```go
package assets

import (
  "bytes"
  "compress/gzip"
  "crypto/md5"
  "encoding/hex"
  "io/ioutil"
  "strings"
)

// returns the contentType for the file
func contentType(filename string) string {
  if strings.HasSuffix(filename, ".png") {
    return "image/png"
  }
  if strings.HasSuffix(filename, ".svg") {
    return "image/png"
  }
  if strings.HasSuffix(filename, ".css") {
    return "text/css"
  }
  if strings.HasSuffix(filename, ".js") {
    return "application/js"
  }
  if strings.HasSuffix(filename, ".eot") {
    return "font/eot"
  }
  if strings.HasSuffix(filename, ".ttf") {
    return "font/ttf"
  }
  if strings.HasSuffix(filename, ".woff") || strings.HasSuffix(filename, ".woff2") {
    return "application/font-woff"
  }
  if strings.HasSuffix(filename, ".html") {
    return "text/html"
  }
  return ""
}

// Asset Gets the file from system if debug otherwise gets it from the stored
// data returns the data, the md5 hash of its content and its content type
// in production the base parameter is ignored
// use it for developement as the relativepath/to/your/public/folder
// ex: ui/myproject/src/public
func Asset(base, path string) ([]byte, string, string) {
  var data []byte
  var err error
  var b bytes.Buffer
  var file string
  if path == "/" {
    file = base + "index.html"
    data, err = ioutil.ReadFile(file)
  } else {
    file = base + path
    data, err = ioutil.ReadFile(file)
  }
  if err != nil {
    file = base + "index.html"
    data, err = ioutil.ReadFile(file)
  }
  if err != nil {
    return []byte("File Not Found " + file), "", "text/html"
  }
  if data != nil {
    w := gzip.NewWriter(&b)
    w.Write(data)
    w.Close()
    data = b.Bytes()
  }
  sum := md5.Sum(data)
  return data, hex.EncodeToString(sum[1:]), contentType(file)
}
```

Go Gophers!

The power is yours!
