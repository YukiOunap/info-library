package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Artist struct {
	Id             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	RelationsURL   string              `json:"relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var artists []Artist

func GetTemp(w http.ResponseWriter, file string) (t *template.Template) {
	t, err := template.ParseFiles(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	return t
}

func Index(w http.ResponseWriter, r *http.Request) {
	t := GetTemp(w, "templates/index.html")

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	t.Execute(w, artists)
}

func ArtistDetail(w http.ResponseWriter, r *http.Request) {
	t := GetTemp(w, "templates/detail.html")

	r.ParseForm()
	id := r.FormValue("artistID")
	url := r.FormValue("artistDetail")

	var targetArtist Artist
	intId, _ := strconv.Atoi(id)
	for _, artist := range artists {
		if artist.Id == intId {
			targetArtist = artist
			break
		}
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&targetArtist)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	t.Execute(w, targetArtist)
}

func main() {
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/detail", ArtistDetail)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
