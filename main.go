package main

// Newest
import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	http1   = false
	http2   = false
	https   = false
	port    = 8000
	rootdir = "view"
	conf    = "config.yaml"
	crt     = "server.crt"
	key     = "server.key"
)

func main() {
	fhttp1 := flag.Bool("http1", false, "server protocol: HTTP/1.1")
	fhttps := flag.Bool("https", false, "server protocol: HTTPS")
	fport := flag.Int("port", 8000, "port number")
	frootdir := flag.String("rootdir", "view", "server root directory")
	fcrt := flag.String("crt", "server.crt", "server certificate (HTTPS)")
	fkey := flag.String("key", "server.key", "server key (HTTPS)")
	fhelp := flag.Bool("help", false, "help")
	flag.Parse()
	fs := http.FileServer(http.Dir(*frootdir))
	http.Handle("/", fs)
	if *fhelp {
		fmt.Println("help")
	} else if *fhttp1 {
		fmt.Println("ready")
		http.ListenAndServe(":"+strconv.Itoa(*fport), nil)
	} else if *fhttps {
		fmt.Println("ready")
		err := http.ListenAndServeTLS(":"+strconv.Itoa(*fport), *fcrt, *fkey, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("ready")
		http.ListenAndServe(":"+strconv.Itoa(port), nil)
	}
}
