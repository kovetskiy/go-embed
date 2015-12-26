package main

import (
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	newline    = []byte{'\n'}
	dataindent = []byte{'\t'}
	space      = []byte{' '}

	input  = flag.String("input", "", "The path to the folder containing the assets")
	output = flag.String("output", "", "The output filename")
	tag    = flag.String("tag", "", "The tag to use for the generated package")

	files       = map[string]string{}
	filesType   = map[string]string{}
	regFuncName = regexp.MustCompile(`[^a-zA-Z0-9_]`)
)

// ByteWriter takes text input and writes them as hexadecimal bytes
type ByteWriter struct {
	io.Writer
	c      int
	digest hash.Hash
	hashed []byte
}

func (w *ByteWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	for n = range p {
		if w.c%12 == 0 {
			w.Writer.Write(newline)
			w.Writer.Write(dataindent)
			w.c = 0
		} else {
			w.Writer.Write(space)
		}

		fmt.Fprintf(w.Writer, "0x%02x,", p[n])
		w.c++
	}

	n++
	w.digest.Write(p)
	w.hashed = w.digest.Sum(w.hashed)
	return
}

func safeFunctionName(name string) string {
	var inBytes, outBytes []byte
	var toUpper bool
	// uppercase the first char to make it public
	name = strings.Title(name)
	inBytes = []byte(name)
	for i := 0; i < len(inBytes); i++ {
		if regFuncName.Match([]byte{inBytes[i]}) {
			toUpper = true
		} else if toUpper {
			outBytes = append(outBytes, []byte(strings.ToUpper(string(inBytes[i])))...)
			toUpper = false
		} else {
			outBytes = append(outBytes, inBytes[i])
		}
	}
	name = string(outBytes)
	return name
}

// need to pipe read to write
func recursiveRead(w io.Writer, folder string) {
	fileInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	for _, info := range fileInfos {
		filename := path.Join(folder, info.Name())
		if info.IsDir() {
			recursiveRead(w, filename)
		} else {
			println("Reading File -> " + filename)
			relativePath := strings.Replace(filename, *input, "", -1)
			relativePath = path.Join("/", relativePath)
			fd, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer fd.Close()
			_, err = fmt.Fprintf(w, "var %s = []byte{", safeFunctionName(relativePath))
			if err != nil {
				panic(err)
			}
			byteWriter := &ByteWriter{
				Writer: w,
				digest: md5.New(),
				hashed: []byte{},
			}
			gz, err := gzip.NewWriterLevel(byteWriter, gzip.BestCompression)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(gz, fd)
			if err != nil {
				panic(err)
			}
			err = gz.Close()
			if err != nil {
				panic(err)
			}
			_, err = fmt.Fprintf(w, "\n}\n")
			if err != nil {
				panic(err)
			}
			files[relativePath] = hex.EncodeToString(byteWriter.digest.Sum(nil))
			filesType[relativePath] = contentType(relativePath)
		}
	}
}

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

func randStr() string {
	dictionary := "0123456789abcdef"
	bytes := make([]byte, 34)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func main() {
	flag.Parse()
	if *input == "" {
		flag.PrintDefaults()
		panic("-input is required.")
	}
	if *output == "" {
		flag.PrintDefaults()
		panic("-output is required.")
	}
	outputFile, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	bfd := bufio.NewWriter(outputFile)
	defer bfd.Flush()
	if _, err = fmt.Fprint(bfd, "// Code generated by go-embed\n\n"); err != nil {
		panic(err)
	}
	if *tag != "" {
		if _, err = fmt.Fprint(bfd, `// +build `+*tag+"\n\n"); err != nil {
			panic(err)
		}
	}
	if _, err = fmt.Fprint(bfd, "package assets\n"); err != nil {
		panic(err)
	}
	if _, err = fmt.Fprint(bfd, `
import (
    "bytes"
    "compress/gzip"
    "crypto/md5"
    "encoding/hex"
    "io/ioutil"
    "strings"
)
	
`); err != nil {
		panic(err)
	}
	recursiveRead(bfd, *input)
	if _, err = fmt.Fprintf(bfd, `
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
func Asset(path string, debug bool) ([]byte, string, string) {
    if debug {
    	var data []byte
      	var err error
      	var b bytes.Buffer
      	var file string
      	if path == "/" {
        	file = "%s" + "index.html"
        	data, err = ioutil.ReadFile(file)
      	} else {
        	file = "%s" + path
        	data, err = ioutil.ReadFile(file)
      	}
      	if err != nil {
        	file = "%s" + "index.html"
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
      	return data, hex.EncodeToString(sum[1:len(sum)]), contentType(file)
    }
    switch path {
`, *input, strings.TrimSuffix(*input, "/"), *input); err != nil {
		panic(err)
	}
	indexFile := `[]byte{}`
	if _, ok := files["/index.html"]; ok {
		indexFile = "IndexHtml"
	}
	for path, hash := range files {
		if _, err = fmt.Fprintf(bfd, `	case "%s":
		return %s, "%s", "%s"
`, path, safeFunctionName(path), hash, filesType[path]); err != nil {
			panic(err)
		}
	}
	if _, err = fmt.Fprintf(bfd, `	default:
    	return %s, "%s", "text/html"
	}
}
`, indexFile, randStr()); err != nil {
		panic(err)
	}
}
