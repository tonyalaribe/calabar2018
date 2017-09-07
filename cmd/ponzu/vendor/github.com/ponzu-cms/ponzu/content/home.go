package content

import (
	"net/http"

	"github.com/bosssauce/frontend"
)

func init() {
	frontend.Router.HandleFunc("/assets/{rest:.*}", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("%#v", r.URL)
		// log.Println("./site/" + r.URL.Path[1:])
		http.ServeFile(w, r, "./site/"+r.URL.Path[1:])
		// http.FileServer(http.Dir("./public/").Open(name)
	})

	frontend.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		http.ServeFile(w, r, "./site/index.html")
	})
}
