package content

import (
	"encoding/json"
	"errors"
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

var BASEURL = "http://localhost:8080"

func RecoverWrap(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				log.Println(err.Error())
				http.Error(w, "404 Page not found", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func init() {
	url := os.Getenv("DOMAIN")
	if url != "" {
		BASEURL = url
	}
	templates = ParseTemplates("./site")
	templates = ParseTemplates("./site/partials")

	frontend.Router.HandleFunc("/assets/{rest:.*}", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("%#v", r.URL)
		// log.Println("./site/" + r.URL.Path[1:])
		http.ServeFile(w, r, "./site/"+r.URL.Path[1:])
		// http.FileServer(http.Dir("./public/").Open(name)
	}))

	frontend.Router.HandleFunc("/hotels", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		data := make(map[string]interface{})
		response, err := http.Get(BASEURL + "/api/contents?type=Hotel")
		if err != nil {
			log.Printf("%s\n", err)
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)

			// json.Unmarshal(data, v)
			json.Unmarshal(body, &data)
			// log.Printf("%s\n", string(body))
		}

		renderTemplate(w, "hotels.html", data)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))

	frontend.Router.HandleFunc("/schedule", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		data := make(map[string]interface{})
		response, err := http.Get(BASEURL + "/api/contents?type=Schedule")
		if err != nil {
			log.Printf("%s\n", err)
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)

			// json.Unmarshal(data, v)
			json.Unmarshal(body, &data)
			// log.Printf("%s\n", string(body))
		}

		renderTemplate(w, "schedule.html", data)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))

	frontend.Router.HandleFunc("/", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		data := make(map[string]interface{})
		response, err := http.Get(BASEURL + "/api/contents?type=Sponsor")
		if err != nil {
			log.Printf("%s\n", err)
		} else {
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Printf("%s\n", err)
			}
			sponsor := make(map[string]interface{})
			// json.Unmarshal(data, v)
			json.Unmarshal(body, &sponsor)
			data["sponsors"] = sponsor["data"]

			// log.Printf("%#v\n", data)
		}

		renderTemplate(w, "index.html", data)
		// http.ServeFile(w, r, "./site/index.html")
	}))

}
