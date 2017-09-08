package content

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bosssauce/frontend"
)

var templates = template.New("").Funcs(template.FuncMap{
	"odd": func(number int) bool {

		if number%2 == 0 {
			return false
		}
		return true
	},
})

// Render a template given a model
func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ParseTemplates(folder string) *template.Template {
	templ := templates
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			log.Println(path)
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		panic(err)
	}

	return templ
}

const BASEURL = "http://localhost:8080"

func init() {
	templates = ParseTemplates("./site")
	templates = ParseTemplates("./site/partials")

	frontend.Router.HandleFunc("/assets/{rest:.*}", func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("%#v", r.URL)
		// log.Println("./site/" + r.URL.Path[1:])
		http.ServeFile(w, r, "./site/"+r.URL.Path[1:])
		// http.FileServer(http.Dir("./public/").Open(name)
	})

	frontend.Router.HandleFunc("/hotels", func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		response, err := http.Get(BASEURL + "/api/contents?type=Hotel")
		if err != nil {
			log.Printf("%s\n", err)
		}
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		data := make(map[string]interface{})
		// json.Unmarshal(data, v)
		json.Unmarshal(body, &data)
		log.Printf("%s\n", string(body))
		renderTemplate(w, "hotels.html", data)
		// http.ServeFile(w, r, "./site/hotels.html")
	})

	frontend.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		http.ServeFile(w, r, "./site/index.html")
	})

}
