package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"strings"
)

func contentType(filename string) string {
	if strings.HasSuffix(filename, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(filename, ".svg") {
		return "image/svg"
	}
	if strings.HasSuffix(filename, ".jpg") {
		return "image/jpg"
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
	return data, hex.EncodeToString(sum[1:]), contentType(file), nil
}
