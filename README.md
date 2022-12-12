# serve

## Summary

Simple tool to allow a directory to be served over HTTP(s), think `python -m http.server` when this is unavilable in the current environment, primarily aimed at providing a host for PoCs. 

Also supports some additional options:

* Support for HTTPS when provided a PEM key & cert
* CORS headers to allow client requests from anywhere to permitted in the browser
* Additional custom headers to be returned in the response

## Releases

Updated versions are automatically built and released, please see https://github.com/stavinski/serve/releases.

## Usage

~~~
USAGE: ./serve [options] <ADDR>

ADDR: Binding address to use, can be just the port (:8000), or the IP/hostname and the port (127.0.0.1:8000) to restrict only localhost. 

OPTIONS:
  -d, --dir             Directory to serve files from, defaults to the cwd
  -s, --secure  Use HTTPS. Requires cert and key pair be provided
  -c, --cert    Certificate file to use in PEM format
  -k, --key             Key file to use in PEM format
  --headers             Add extra header(s). Expected to be in name:value format and comma separated
  --cors                Add CORS header to allow calls from any origin (Access-Control-Allow-Origin: *)

EXAMPLES:
  ./serve -d public :8000
        Serve files over HTTP on any IP over port 8000 from the public directory
  ./serve -s -c cert.pem -k key.pem 127.0.0.1:443
        Serve files from cwd over for localhost only using HTTPS on port 443
  ./serve --headers 'X-Foo: Test' 192.168.1.10:8000
        Serve files from cwd on 192.168.1.10 over port 8000 with extra X-Foo HTTP header in response
~~~
