package main

// Newest
import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/net/http2"
	"gopkg.in/yaml.v2"
)

type conf struct {
	Port                 int
	Http, Dir, Cert, Key string
}

var dcnf = conf{
	Http: "1",
	Port: 8000,
	Dir:  "dir_gokt",
	Cert: "crt_gokt.crt",
	Key:  "key_gokt.key",
}
var dconf = "cnf_gokt.yaml"

func (c *conf) readconf(i string) *conf {
	yamlFile, err := ioutil.ReadFile(i)
	if err != nil && i != "" {
		fmt.Printf("error: %v file not found, switch using default configuration\n", i)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("error: can't read %v file, switch using default configuration\n", i)
	}
	return c
}

func help() {
	fmt.Println("")
	fmt.Println("  < Go-kart Track Web Server v1.0 >")
	fmt.Println("")
	fmt.Println("    This web server will use parameter input, configuration file (.yaml), or hard-coded configuration (if")
	fmt.Println("no parameter input or configuration file given) to run. To run with parameter input, use the followings:")
	fmt.Println("    > '-conf' : The configuration file used by this server (.yaml)")
	fmt.Println("    > '-http' : The protocol used by this server ('1' for HTTP/1.1 and 's'/'S' for HTTPS/2)")
	fmt.Println("    > '-port' : The port number used by this server")
	fmt.Println("    > '-dir'  : The server root directory used by this server")
	fmt.Println("    > '-cert' : The certificate file used by this server (HTTPS protocol)")
	fmt.Println("    > '-key'  : The key file used by this server (HTTPS protocol)")
	fmt.Println("For example: './gokt -http s -port 8080 -conf myConf/conf.yaml', this will tell the server to run using")
	fmt.Println("HTTPS/2 protocol, in port 8080, and read the conf.yaml file in myConf folder for the other configurations.")
	fmt.Println("The inside of the configuration file should mention 'http:', 'port:', 'dir:', 'cert:', 'key:' and their")
	fmt.Println("value. The default configuration file is 'cnf_gokt.yaml' file, default root directory is the folder 'dir_gokt'.")
	fmt.Println("")
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
	if lcnf.Dir == "" || lcnf.Dir == "." || lcnf.Dir == "/" {
		if d.Dir == "" || d.Dir == "." || d.Dir == "/" {
			lcnf.Dir = dcnf.Dir
		} else {
			lcnf.Dir = d.Dir
		}
	}
	fs := http.FileServer(http.Dir(lcnf.Dir))
	http.Handle("/", fs)
	fmt.Println("Go-kart running on HTTP/1.1 |", "server directory:", lcnf.Dir, "| port:", lcnf.Port)
	http.ListenAndServe(":"+strconv.Itoa(lcnf.Port), nil)
}

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
	if lcnf.Dir == "" || lcnf.Dir == "." || lcnf.Dir == "/" {
		if d.Dir == "" || d.Dir == "." || d.Dir == "/" {
			lcnf.Dir = dcnf.Dir
		} else {
			lcnf.Dir = d.Dir
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
	fmt.Println("Go-kart running on HTTPS/2 |", "server directory:", lcnf.Dir, "| port:", lcnf.Port, "| certificate:", lcnf.Cert,
		"| key:", lcnf.Key)
	http.ListenAndServeTLS(":"+strconv.Itoa(lcnf.Port), lcnf.Cert, lcnf.Key, nil)
	fmt.Printf("error: unable to find %v and/or %v file to run on HTTPS/2", lcnf.Cert, lcnf.Key)
}

func servehttp2() {
	var srv http.Server
	srv.Addr = ":8000"
	http2.ConfigureServer(&srv, nil)
	fs := http.FileServer(http.Dir("dir_gokt"))
	http.Handle("/", fs)
	srv.ListenAndServe()

}

func main() {
	var rcnf conf
	var bcnf *conf
	fhelp := flag.Bool("help", false, "Show brief information on how this server works")
	fconf := flag.String("conf", "", "The configuration file used by this server (.yaml)")
	fhttp := flag.String("http", "", "The protocol used by this server ('1' for HTTP/1.1 and 's'/'S' for HTTPS/2)")
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
		} else {
			fmt.Printf("error: unknown protocol %v and %v file not found, switch using default configuration\n", *fhttp, *fconf)
			bcnf = rcnf.readconf(dconf)
			if bcnf.Http == "1" {
				servehttp1(dconf, fcnf)
			} else if bcnf.Http == "s" || bcnf.Http == "S" {
				servehttps(dconf, fcnf)
			} else {
				servehttp1("", fcnf)
			}
		}
	} else if *fhttp == "" && *fconf != "" {
		bcnf = rcnf.readconf(*fconf)
		if bcnf.Http == "1" {
			servehttp1(*fconf, fcnf)
		} else if bcnf.Http == "s" || bcnf.Http == "S" {
			servehttps(*fconf, fcnf)
		} else {
			bcnf = rcnf.readconf(dconf)
			if bcnf.Http == "1" {
				servehttp1(dconf, fcnf)
			} else if bcnf.Http == "s" || bcnf.Http == "S" {
				servehttps(dconf, fcnf)
			} else {
				servehttp1("", fcnf)
			}
		}
	} else if *fhttp != "" && *fconf == "" {
		if *fhttp == "1" {
			servehttp1(dconf, fcnf)
		} else if *fhttp == "s" || *fhttp == "S" {
			servehttps(dconf, fcnf)
		} else if *fhttp == "2" {
			servehttp2()
		} else {
			fmt.Printf("error: unknown protocol %v, switch using default configuration\n", *fhttp)
			bcnf = rcnf.readconf(dconf)
			if bcnf.Http == "1" {
				servehttp1(dconf, fcnf)
			} else if bcnf.Http == "s" || bcnf.Http == "S" {
				servehttps(dconf, fcnf)
			} else {
				servehttp1("", fcnf)
			}
		}
	} else if *fhttp == "" && *fconf == "" {
		bcnf = rcnf.readconf(dconf)
		if bcnf.Http == "1" {
			servehttp1(dconf, fcnf)
		} else if bcnf.Http == "s" || bcnf.Http == "S" {
			servehttps(dconf, fcnf)
		} else {
			servehttp1("", fcnf)
		}
	}
}
