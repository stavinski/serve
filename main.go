// A basic HTTP(S) server to be used for checking proof of concepts, think when 'python -m http.server' not available.
//
// Also supports:
//
// * HTTPS with a PEM key/cert
//
// * CORS header for allowing client requests from anywhere
//
// * Additional custom headers in the response
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var version = "x.x.x" // overwritten by build

// Options set for the program
type Options struct {
	addr     string
	dir      string
	certFile string
	keyFile  string
	useHTTPS bool
	useCORS  bool
	headers  map[string]string
}

func handleFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	var usage = `USAGE: ` + os.Args[0] + ` (v` + version + `) [options] <ADDR>

ADDR: Binding address to use, can be just the port (:8000), or the IP/hostname and the port (127.0.0.1:8000) to restrict only localhost. 
	
OPTIONS:
  -d, --dir		Directory to serve files from, defaults to the cwd
  -s, --secure	Use HTTPS. Requires cert and key pair be provided
  -c, --cert	Certificate file to use in PEM format
  -k, --key		Key file to use in PEM format
  --headers		Add extra header(s). Expected to be in name:value format and comma separated
  --cors		Add CORS header to allow calls from any origin (Access-Control-Allow-Origin: *)

EXAMPLES:
  ` + os.Args[0] + ` -d public :8000							
  	Serve files over HTTP on any IP over port 8000 from the public directory
  ` + os.Args[0] + ` -s -c cert.pem -k key.pem 127.0.0.1:443	
  	Serve files from cwd over for localhost only using HTTPS on port 443
  ` + os.Args[0] + ` --headers 'X-Foo: Test' 192.168.1.10:8000			
  	Serve files from cwd on 192.168.1.10 over port 8000 with extra X-Foo HTTP header in response
`
	fmt.Fprint(flag.CommandLine.Output(), usage)
	os.Exit(1)
}

func parseArgs() Options {
	var (
		ret     Options
		headers string
	)
	ret = Options{}

	cwd, err := os.Getwd()
	handleFatalErr(err)

	flag.Usage = usage
	flag.BoolVar(&ret.useHTTPS, "s", false, "")
	flag.BoolVar(&ret.useHTTPS, "secure", false, "")
	flag.BoolVar(&ret.useCORS, "cors", false, "")
	flag.StringVar(&ret.certFile, "c", "", "")
	flag.StringVar(&ret.certFile, "cert", "", "")
	flag.StringVar(&ret.keyFile, "k", "", "")
	flag.StringVar(&ret.keyFile, "key", "", "")
	flag.StringVar(&headers, "headers", "", "")
	flag.StringVar(&ret.dir, "d", cwd, "")
	flag.StringVar(&ret.dir, "dir", cwd, "")
	flag.Parse()

	// addr not been provided
	if len(flag.Args()) < 1 {
		log.Println("ADDR has not been provided.")
		usage()
	}

	// if HTTPS flag set then make sure HTTPS options have also been provided
	if ret.useHTTPS && (len(ret.certFile) == 0 || len(ret.keyFile) == 0) {
		log.Println("Missing cert or key when using secure flag.")
		usage()
	}

	// parse headers into string map
	if len(headers) > 0 {
		ret.headers = make(map[string]string, len(headers))
		vals := strings.Split(headers, ",")
		for _, val := range vals {
			k, v, found := strings.Cut(val, ":")
			if !found {
				log.Println("Invalid headers string was provided.")
				usage()
			}
			ret.headers[k] = v
		}
	}

	ret.addr = flag.Arg(0)
	return ret
}

func getHandler(opts Options) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		defer next(w, req)

		// Add CORS allow header
		if opts.useCORS {
			w.Header().Add("Access-Control-Allow-Origin", "*")
		}

		// Add extra headers to the response
		for key, val := range opts.headers {
			w.Header().Add(key, val)
		}
	}
}

func main() {
	opts := parseArgs()
	r := mux.NewRouter()
	n := negroni.New()
	n.UseFunc(getHandler(opts))
	n.Use(negroni.NewLogger())
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(opts.dir)))
	n.UseHandler(r)
	log.Printf("Serving files from: %s", opts.dir)
	if opts.useHTTPS {
		log.Printf("HTTPS: %s\n", opts.addr)
		err := http.ListenAndServeTLS(opts.addr, opts.certFile, opts.keyFile, n)
		handleFatalErr(err)
	} else {
		log.Printf("HTTP: %s\n", opts.addr)
		err := http.ListenAndServe(opts.addr, n)
		handleFatalErr(err)
	}
}
