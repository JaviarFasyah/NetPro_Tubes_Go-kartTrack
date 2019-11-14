package main

// Newest
import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gopkg.in/yaml.v2"
)

type conf struct {
	Port                 int
	Http, Dir, Cert, Key string
}

var dcnf = conf{
	Http: "1",
	Port: 8000,
	Dir:  "view",
	Cert: "server.crt",
	Key:  "server.key",
}
var dconf = "gokt.yaml"

func (c *conf) readconf(i string) *conf {
	yamlFile, err := ioutil.ReadFile(i)
	if err != nil {
		fmt.Printf("error: %v file not found, switch using default configuration\n", i)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("error: can't read %v file, switch using default configuration\n", i)
	}
	return c
}

func help() {
	fmt.Println("help")
}

func servehttp1(i string, j conf) {
	var c conf
	var d *conf
	var lcnf = conf{
		Port: j.Port,
		Dir:  j.Dir,
	}
	d = c.readconf(i)
	if lcnf.Port == 0 {
		if d.Port != 0 {
			lcnf.Port = d.Port
		} else {
			lcnf.Port = dcnf.Port
		}
	}
	if lcnf.Dir == "" {
		if d.Dir != "" {
			lcnf.Dir = d.Dir
		} else {
			lcnf.Dir = dcnf.Dir
		}
	}
	fs := http.FileServer(http.Dir(lcnf.Dir))
	http.Handle("/", fs)
	fmt.Println("running on HTTP/1.1 |", "server directory:", lcnf.Dir, "| port:", lcnf.Port)
	http.ListenAndServe(":"+strconv.Itoa(lcnf.Port), nil)
}

func servehttp2(i string) {}

func servehttps(i string, j conf) {
	var c conf
	var d *conf
	var lcnf = conf{
		Port: j.Port,
		Dir:  j.Dir,
		Cert: j.Cert,
		Key:  j.Key,
	}
	d = c.readconf(i)
	if lcnf.Port == 0 {
		if d.Port != 0 {
			lcnf.Port = d.Port
		} else {
			lcnf.Port = dcnf.Port
		}
	}
	if lcnf.Dir == "" {
		if d.Dir != "" {
			lcnf.Dir = d.Dir
		} else {
			lcnf.Dir = dcnf.Dir
		}
	}
	if lcnf.Cert == "" {
		if d.Cert != "" {
			lcnf.Cert = d.Cert
		} else {
			lcnf.Cert = dcnf.Cert
		}
	}
	if lcnf.Key == "" {
		if d.Key != "" {
			lcnf.Key = d.Key
		} else {
			lcnf.Key = dcnf.Key
		}
	}
	fs := http.FileServer(http.Dir(lcnf.Dir))
	http.Handle("/", fs)
	fmt.Println("running on HTTPS |", "server directory:", lcnf.Dir, "| port:", lcnf.Port)
	http.ListenAndServeTLS(":"+strconv.Itoa(lcnf.Port), lcnf.Cert, lcnf.Key, nil)
}

func main() {
	var rcnf conf
	fhelp := flag.Bool("help", false, "help")
	fconf := flag.String("conf", "", "The configuration file used by this server (.yaml)")
	fhttp := flag.String("http", "", "The protocol used by this server ('1' for HTTP/1.1, '2' for HTTP/2, 's' for HTTPS)")
	fport := flag.Int("port", 0, "The port number used by this server")
	fdir := flag.String("dir", "", "The server root directory used by this server")
	fcert := flag.String("cert", "", "The certificate file used by this server (HTTPS protocol)")
	fkey := flag.String("key", "", "The key file used by this server (HTTPS protocol)")
	flag.Parse()
	var fcnf = conf{
		Port: *fport,
		Dir:  *fdir,
		Cert: *fcert,
		Key:  *fkey,
	}
	if *fhelp {
		help()
	} else if *fhttp != "" && *fconf != "" {
		if *fhttp == "1" {
			servehttp1(*fconf, fcnf)
		} else if *fhttp == "s" || *fhttp == "S" {
			servehttps(*fconf, fcnf)
		} else if *fhttp == "2" {
			fmt.Println("http2")
		} else {
			servehttp1(dconf, fcnf)
		}
	} else if *fhttp == "" && *fconf != "" {
		if rcnf.readconf(*fconf).Http == "1" {
			servehttp1(*fconf, fcnf)
		} else if rcnf.readconf(*fconf).Http == "s" || rcnf.readconf(*fconf).Http == "S" {
			servehttps(*fconf, fcnf)
		} else if rcnf.readconf(*fconf).Http == "2" {
			fmt.Println("http2")
		} else {
			servehttp1(dconf, fcnf)
		}
	} else if *fhttp != "" && *fconf == "" {
		if *fhttp == "1" {
			servehttp1(dconf, fcnf)
		} else if *fhttp == "s" || *fhttp == "S" {
			servehttps(dconf, fcnf)
		} else if *fhttp == "2" {
			fmt.Println("http2")
		} else {
			servehttp1(dconf, fcnf)
		}
	} else if *fhttp == "" && *fconf == "" {
		servehttp1(dconf, fcnf)
	}
}
