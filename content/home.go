package content

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	"add": func(a, b int) string {
		return strconv.Itoa(a + b)
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
		// log.Println("calling next handler...")
		h(w, r)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	frontend.Router.HandleFunc("/register", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)

		renderTemplate(w, "register.html", nil)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))
	frontend.Router.HandleFunc("/banquet", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)

		renderTemplate(w, "banquet.html", nil)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))

	frontend.Router.HandleFunc("/view/individuals", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		individuals := make(map[string][]RegisteredIndividual)
		Resp, err := http.Get(BASEURL + "/api/contents?type=RegisteredIndividuals")
		if err != nil {
			renderTemplate(w, "viewindividuals.html", nil)
			log.Printf("%s\n", err)
			return
		}

		defer Resp.Body.Close()
		body, err := ioutil.ReadAll(Resp.Body)
		if err != nil {
			log.Println("error: ", err)
		}
		json.Unmarshal(body, &individuals)
		log.Printf("individuals: %#v\nbody: %#v\n", individuals, string(body))
		renderTemplate(w, "viewindividuals.html", individuals)

		// http.ServeFile(w, r, "./site/hotels.html")
	}))
	frontend.Router.HandleFunc("/view/clubs", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		clubs := make(map[string][]RegisteredClub)
		Resp, err := http.Get(BASEURL + "/api/contents?type=RegisteredClubs")
		if err != nil {
			renderTemplate(w, "viewclubs.html", nil)
			log.Printf("%s\n", err)
			return
		}

		defer Resp.Body.Close()
		body, err := ioutil.ReadAll(Resp.Body)
		if err != nil {
			log.Println("error: ", err)
		}
		json.Unmarshal(body, &clubs)
		log.Printf("clubs: %#v\nbody: %#v\n", clubs, string(body))
		renderTemplate(w, "viewclubs.html", clubs)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))
	frontend.Router.HandleFunc("/view/banquets", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		banquets := make(map[string][]Banquet)
		Resp, err := http.Get(BASEURL + "/api/contents?type=Banquet")
		if err != nil {
			renderTemplate(w, "viewbanquets.html", nil)
			log.Printf("%s\n", err)
			return
		}

		defer Resp.Body.Close()
		body, err := ioutil.ReadAll(Resp.Body)
		if err != nil {
			log.Println("error: ", err)
		}
		json.Unmarshal(body, &banquets)
		log.Printf("Banquets: %#v\nbody: %#v\n", banquets, string(body))
		renderTemplate(w, "viewbanquets.html", banquets)
	}))

	frontend.Router.HandleFunc("/add_newsletter_email", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&data)
		log.Println("data: ", data)
		var buf bytes.Buffer
		ww := multipart.NewWriter(&buf)

		log.Println("writing...")

		ww.WriteField("email", data["Email"].(string))
		ww.Close()

		newsletters := make(map[string][]Newsletter)
		Resp, err := http.Get(BASEURL + "/api/contents?type=Newsletter")
		if err != nil {
			log.Printf("%s\n", err)
			return
		} else {
			defer Resp.Body.Close()
			body, err := ioutil.ReadAll(Resp.Body)
			if err != nil {
				log.Println("error: ", err)
			}
			json.Unmarshal(body, &newsletters)
			log.Printf("%#v\n", newsletters)
			for _, newsletter := range newsletters["data"] {
				if newsletter.Email == data["Email"].(string) {
					json.NewEncoder(w).Encode(map[string]string{"error": "Already Subscribed"})
					return
				}
			}
		}

		log.Println("Making a post request...")
		resp, err := http.Post(BASEURL+"/api/content/create?type=Newsletter", ww.FormDataContentType(), &buf)
		if err != nil {
			log.Println(err)
		}
		byt, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(byt))
		json.NewEncoder(w).Encode(data)

	}))

	frontend.Router.HandleFunc("/register_banquet", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&data)
		log.Println("data: ", data)
		var buf bytes.Buffer
		ww := multipart.NewWriter(&buf)

		log.Println("writing...")

		err := ww.WriteField("registration_id", data["RegistrationID"].(string))
		log.Println(err)
		err = ww.WriteField("phone", data["Phone"].(string))
		log.Println(err)
		err = ww.WriteField("email", data["Email"].(string))
		log.Println(err)
		log.Printf("type: %T", data["Amount"])
		err = ww.WriteField("amount", strconv.FormatFloat(data["Amount"].(float64), 'f', -1, 64))
		log.Println("err: ", err)

		ww.Close()

		// check if registration id is present...
		banquets := make(map[string][]Banquet)
		BanquetResp, err := http.Get(BASEURL + "/api/contents?type=Banquet")
		if err != nil {
			log.Printf("%s\n", err)
			return
		} else {
			defer BanquetResp.Body.Close()
			body, _ := ioutil.ReadAll(BanquetResp.Body)

			json.Unmarshal(body, &banquets)
			log.Printf("%#v\n", banquets)
			for _, banq := range banquets["data"] {
				if banq.RegistrationId == data["RegistrationID"].(string) {
					json.NewEncoder(w).Encode(map[string]string{"error": "Already Registered"})
					return
				}
			}
		}

		isRegistered := false
		users := make(map[string][]RegisteredIndividual)
		UsersResp, err := http.Get(BASEURL + "/api/contents?type=RegisteredIndividual")
		if err != nil {
			log.Printf("%s\n", err)
			return
		} else {
			defer UsersResp.Body.Close()
			body, _ := ioutil.ReadAll(UsersResp.Body)

			json.Unmarshal(body, &users)
			log.Printf("%#v\n", users)
			for _, user := range users["data"] {
				if user.RegisterID == data["RegistrationID"].(string) {
					isRegistered = true
					break
				}
			}
		}
		if !isRegistered {
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Registration ID"})
			return
		}

		log.Println("Making a post request....")
		resp, err := http.Post(BASEURL+"/api/content/create?type=Banquet", ww.FormDataContentType(), &buf)
		if err != nil {
			log.Println(err)
		}
		byt, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(byt))
		json.NewEncoder(w).Encode(data)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))

	frontend.Router.HandleFunc("/register_individual", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)

		data := make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&data)
		log.Println(data)
		var buf bytes.Buffer
		ww := multipart.NewWriter(&buf)
		ww.WriteField("name", data["FullName"].(string))
		ww.WriteField("club", data["Club"].(string))
		ww.WriteField("region", data["Region"].(string))
		ww.WriteField("district", data["District"].(string))
		ww.WriteField("phone", data["Phone"].(string))
		ww.WriteField("email", data["Email"].(string))
		ww.Close()

		_, err := http.Post(BASEURL+"/api/content/create?type=RegisteredIndividuals", ww.FormDataContentType(), &buf)
		if err != nil {
			log.Println(err)
		}
		// byt, _ := ioutil.ReadAll(resp.Body)
		json.NewEncoder(w).Encode(data)
		// renderTemplate(w, "register.html", nil)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))

	frontend.Router.HandleFunc("/register_club", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)

		data := make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&data)
		log.Println(data)
		var buf bytes.Buffer
		ww := multipart.NewWriter(&buf)
		ww.WriteField("name", data["FullName"].(string))
		ww.WriteField("club", data["Club"].(string))
		ww.WriteField("region", data["Region"].(string))
		ww.WriteField("district", data["District"].(string))
		ww.WriteField("phone", data["Phone"].(string))
		ww.WriteField("email", data["Email"].(string))
		ww.Close()

		_, err := http.Post(BASEURL+"/api/content/create?type=RegisteredClubs", ww.FormDataContentType(), &buf)
		if err != nil {
			log.Println(err)
		}
		// byt, _ := ioutil.ReadAll(resp.Body)
		json.NewEncoder(w).Encode(data)
		// renderTemplate(w, "register.html", nil)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))

	frontend.Router.HandleFunc("/register_booking", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)

		data := make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&data)
		var buf bytes.Buffer
		ww := multipart.NewWriter(&buf)
		ww.WriteField("full_name", data["FullName"].(string))
		ww.WriteField("phone", data["Phone"].(string))
		ww.WriteField("email", data["Email"].(string))
		ww.WriteField("room", data["Room"].(string))
		ww.WriteField("hotel", data["Hotel"].(string))
		ww.Close()

		_, err := http.Post(BASEURL+"/api/content/create?type=Bookings", ww.FormDataContentType(), &buf)
		if err != nil {
			log.Println(err)
		}
		// byt, _ := ioutil.ReadAll(resp.Body)
		json.NewEncoder(w).Encode(data)
		// renderTemplate(w, "register.html", nil)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))

	frontend.Router.HandleFunc("/hotels", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		hotels := make(map[string][]Hotel)
		response, err := http.Get(BASEURL + "/api/contents?type=Hotel")
		if err != nil {
			log.Printf("%s\n", err)
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)

			// json.Unmarshal(hotels, v)
			json.Unmarshal(body, &hotels)
			// log.Printf("%s\n", string(body))
		}

		rooms := make(map[string][]Room)
		response, err = http.Get(BASEURL + "/api/contents?type=Room")
		if err != nil {
			log.Printf("%s\n", err)
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)

			// json.Unmarshal(hotels, v)
			json.Unmarshal(body, &rooms)
		}

		finalHotels := []Hotel{}
		for _, hotel := range hotels["data"] {
			log.Println(hotel)
			nrooms := []Room{}
			for _, room := range rooms["data"] {

				if room.Hotel == fmt.Sprintf("/api/content?type=Hotel&id=%d", hotel.ID) {
					log.Println("match")
					nrooms = append(nrooms, room)
				}
			}
			hotel.Rooms = nrooms
			finalHotels = append(finalHotels, hotel)
		}

		data := make(map[string][]Hotel)
		data["data"] = finalHotels
		renderTemplate(w, "hotels.html", data)
		// http.ServeFile(w, r, "./site/hotels.html")
	}))

	frontend.Router.HandleFunc("/hotels/book", RecoverWrap(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		log.Println(r.URL.Query())

		// vars := mux.Vars(r)
		// log.Printf("%#v", vars)
		// slug := vars["slug"]
		// log.Println(slug)

		// hotelSlug := r.URL.Query().Get("hotel")
		roomSlug := r.URL.Query().Get("room")

		room := make(map[string][]Room)
		response, err := http.Get(BASEURL + "/api/content?type=Room&slug=" + roomSlug)
		if err != nil {
			log.Printf("%s\n", err)
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			// log.Printf("> %#v", string(body))
			// json.Unmarshal(hotels, v)
			json.Unmarshal(body, &room)
			// log.Printf("%#v\n", room)
			// log.Printf("%s\n", string(body))
		}

		hotel := make(map[string][]Hotel)
		response, err = http.Get(BASEURL + room["data"][0].Hotel)
		if err != nil {
			log.Printf("%s\n", err)
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)

			// json.Unmarshal(hotels, v)
			json.Unmarshal(body, &hotel)
			// log.Printf("%#v\n", hotel)
		}

		data := make(map[string]interface{})
		if len(hotel["data"]) > 0 {
			data["hotel"] = hotel["data"][0]
		}
		if len(room["data"]) > 0 {
			data["room"] = room["data"][0]
		}
		// log.Println(data)
		renderTemplate(w, "book_hotel.html", data)
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
