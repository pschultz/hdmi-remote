package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	var dev bool
	var bundle bool

	flag.BoolVar(&dev, "dev", false, "serve index.html from disk")
	flag.BoolVar(&bundle, "bundle", false, "re-generate index.html.go")
	flag.Parse()

	if bundle {
		makeIndexHTML()
		return
	}

	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/switch/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}

		port, _ := strconv.Atoi(r.URL.Path[8:])
		if port < 1 || port > 4 {
			http.Error(w, "Unparsable or invalid port", 422)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		switch err := changePort(ctx, port); err {
		case nil:
			w.WriteHeader(200)
		case context.Canceled:
			http.Error(w, err.Error(), 504)
		default:
			http.Error(w, err.Error(), 500)
		}
	})

	fs := http.FileServer(http.Dir("."))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		log.Println(r.Context().Deadline())

		if dev {
			fs.ServeHTTP(w, r)
		} else {
			w.Header().Set("content-type", "text/html; charset=utf-8")
			w.Header().Set("content-encoding", "gzip")
			w.Write(indexHTML)
		}
	})

	srv := &http.Server{
		Addr: ":http",

		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
	}
	if port := os.Getenv("PORT"); port != "" {
		srv.Addr = ":" + port
	}

	log.Println("Starting server on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// sem ensures that we don't attempt to send mulitiple concurrent IR signals.
var sem = make(chan struct{}, 1)

// changePort changes the active input port on the HDMI switch. The context
// must have a deadline.
func changePort(ctx context.Context, port int) error {
	select {
	case sem <- struct{}{}:
		defer func() { <-sem }()
	case <-ctx.Done():
		log.Println("timeout waiting for lock")
		return ctx.Err()
	}

	// Non-empirical tests show the switch takes a few milliseconds to change
	// the port. Bail if there's not enough time left.
	dl, _ := ctx.Deadline()
	if d := time.Until(dl); d < 100*time.Millisecond {
		log.Println("not enough time left after grabbing lock:", d)
		return context.Canceled // not true, but simplifies error handling in caller
	}

	var cmd uint8
	switch port {
	case 1:
		cmd = 0xAA // Button1: 1010 1010b (transmission order)
	case 2:
		cmd = 0xA8 // Button2: 1010 1000b (transmission order)
	case 3:
		cmd = 0xBA // Button3: 1011 1010b (transmission order)
	case 4:
		cmd = 0x38 // Button4: 0011 1000b (transmission order)
	}

	data := uint32(0x00ff0000) // address is always zero
	data |= uint32(cmd) << 8
	data |= uint32(^cmd)

	log.Println("Switching to port", port)
	return exec.CommandContext(ctx, "irsling", fmt.Sprintf("%032b", data)).Run()
}

func makeIndexHTML() {
	b, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}

	buf := &bytes.Buffer{}
	gw := gzip.NewWriter(buf)
	if _, err = gw.Write(b); err != nil {
		log.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("index.html.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "package main")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "// Generated from index.html. Do not edit.")
	fmt.Fprintln(f, "")
	fmt.Fprintf(f, "var indexHTML = %#v\n", buf.Bytes())

	fmt.Println("index.html.go re-generated. Don't forget to run `go build`.")

	return
}
