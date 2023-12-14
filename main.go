package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	t.Execute(w, artists)
}

func ArtistDetail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t, err := template.ParseFiles("templates/detail.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	ID := r.FormValue("artistID")
	fmt.Println("ID:", ID)

	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	var targetArtist Artist
	intID, _ := strconv.Atoi(ID)
	for _, artist := range artists {
		if artist.Id == intID {
			targetArtist = artist
			break
		}
	}

	URL := r.FormValue("artistDetail")
	fmt.Println("URL:", URL)

	response, err = http.Get(URL)
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

	http.ListenAndServe(":8080", nil)
}
