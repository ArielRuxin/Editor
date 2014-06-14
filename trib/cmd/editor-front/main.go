// Front end launcher program.
package main

import (
	"encoding/json"
	"flag"
	"go/build"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"trib"
	"trib/randaddr"
	"trib/ref"
	"triblab"
)

var (
	verbose = flag.Bool("v", false, "verbose logging")
	lab     = flag.Bool("lab", false, "use lab implementation")
	addr    = flag.String("addr", "localhost:rand", "serve address")
	frc     = flag.String("rc", trib.DefaultRCPath,
		"bin storage config file")
	dbinit = flag.Bool("init", false, "do not populate with test data")

	server trib.Server
)

func handleApi(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/")

	reply := func(obj interface{}) {
		bytes, e := json.Marshal(obj)
		noError(e)

		_, e = w.Write(bytes)
		logError(e)
	}

	bytes, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(e)
		return
	}
	input := string(bytes)

	if !*verbose {
		log.Println(name, input)
	}

	var _cuser, _cfile;

	switch name {
	case "add-user":
		e = server.Hello()
		if e != nil {
			reply(NewBool(e == nil, e))
			break
		}
		_cuser = input;
		reply(NewBool(e == nil, e))

	case "search-file":
		_cfile = input
		ret, content, e = server.SearchFile(_cuser, _cfile)
		if e != nil {
			
		}
		reply(NewFile(ret, content, e))

	case "update-file":
		edition := new(Log)
		e := json.Unmarshal(bytes, edition)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		e = server.UpdateFile(_cuser, log)
		reply(NewBool(e == nil, e))

	case "latest":
		/*p := new(Post)
		e := json.Unmarshal(bytes, p)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		e = server.Post(p.Who, p.Message, p.Clock)*/
		reply(NewBool(e == nil, e))

	default:
		w.WriteHeader(404)
	}
}

func ne(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func makeServer() trib.Server {
	if !*lab {
		return ref.NewServer()
	}

	rc, e := trib.LoadRC(*frc)
	ne(e)

	c := triblab.NewBinClient(rc.Backs)

	return triblab.NewFront(c)
}

func wwwPath() string {
	pkg, e := build.Import("trib", "./", build.FindOnly)
	if e != nil {
		log.Fatal(e)
	}
	return filepath.Join(pkg.Dir, "www")
}

func main() {
	flag.Parse()

	server = makeServer()
	if *dbinit {
	//	populate(server)
	}
	*addr = randaddr.Resolve(*addr)
	log.Printf("serve on %s", *addr)

	http.Handle("/", http.FileServer(http.Dir(wwwPath())))
	http.HandleFunc("/api/", handleApi)

	for {
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
