package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"plugin"
	"strings"
	"sync"
)

const addr = ":8080"

type nominator struct {
	sync.RWMutex
	name string
}

func (n *nominator) set(name string) {
	n.Lock()
	defer n.Unlock()
	n.name = name
}

func (n *nominator) get() string {
	n.RLock()
	defer n.RUnlock()
	return n.name
}

var (
	n     = &nominator{}
	cache = make(map[string]func() string) // func-level cache
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", getWho())
	})

	http.HandleFunc("/nominator", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		switch method {
		case "GET":
			fmt.Fprintf(w, "%s", n.get())
		case "PUT":
			if r.Body == nil {
				fmt.Fprintf(w, "no body")
				return
			}
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
				return
			}
			n.set(strings.TrimSpace(string(b)))
			fmt.Fprintf(w, "nominator set to %q", n.get())
		default:
			fmt.Fprintf(w, "unknown method %q (only GET, PUT are in use)", method)
		}
	})

	log.Printf("Server is up at %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getWho() string {
	pn := n.get()

	fKey := fmt.Sprintf("%s#%s", pn, "Who")
	if f, ok := cache[fKey]; ok {
		return f()
	}

	pluginName := fmt.Sprintf(pn + ".so")
	p, err := plugin.Open(pluginName)
	if err != nil {
		return err.Error()
	}

	v, err := p.Lookup("Who")
	if err != nil {
		return err.Error()
	}

	if f, ok := v.(func() string); ok {
		cache[fKey] = f
		return f()
	}
	return fmt.Sprintf("Who() is not %q in %q", "func() string", pluginName)
}
