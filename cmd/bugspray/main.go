package main

import (
	"flag"
	"net/http"
	"net/url"

	"github.com/mcpherrinm/bugspray/internal/proxy"
	"github.com/mcpherrinm/bugspray/internal/ui"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	listen := flag.String("listen", ":8080", "Listen Address")
	target := flag.String("target", "http://localhost:8081", "Target URL")
	flag.Parse()

	acmeProxy := proxy.New(must(url.Parse(*target)))
	ui := ui.Assets()

	mux := http.NewServeMux()
	mux.Handle("/ui/", http.StripPrefix("/ui/", ui))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect / to the /ui/, and proxy everything else
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/ui/", http.StatusMovedPermanently)
			return
		}

		acmeProxy.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(*listen, mux)
	if err != nil {
		panic(err)
	}
}
