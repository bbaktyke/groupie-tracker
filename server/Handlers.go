package server

import (
	"net/http"
	"strconv"
	"text/template"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("htmlfiles/index.html")
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}
	if r.URL.Path != "/" {
		err := "404\nPage not found"
		ErrorPage(w, err, 404)
		return
	}
	if r.Method != http.MethodGet {
		err := "405 method not allowed"
		ErrorPage(w, err, 405)
		return
	}
	artists := MainParser()
	location := GetAllLocations()
	allInfo := Everything{artists, location}

	err = templates.Execute(w, allInfo)
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}
}

func ArtistPage(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("htmlfiles/artists.html")
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}
	if r.URL.Path != "/artists/" {
		err := "404\nPage not found"
		ErrorPage(w, err, 404)
		return
	}
	if r.Method != http.MethodGet {
		err := "405 method not allowed"
		ErrorPage(w, err, 405)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id <= 0 || id >= 53 {
		err := "404\nPage not found"
		ErrorPage(w, err, 404)
		return
	}
	var f General = General{
		Artist:         idParser(id),
		RelationStruct: RelationParser(id),
	}

	err = templates.Execute(w, f)
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}
}

func SearchHadler(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("htmlfiles/searched.html")
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}
	if r.URL.Path != "/search/" {
		err := "404\nPage not found"
		ErrorPage(w, err, 404)
		return
	}
	if r.Method != http.MethodGet {
		err := "405 method not allowed"
		ErrorPage(w, err, 405)
		return
	}
	searched := r.FormValue("Search")
	x := Finder(searched)
	err = templates.Execute(w, x)
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}
}

func FilterHadler(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("htmlfiles/searched.html")
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}
	if r.URL.Path != "/filter/" {
		err := "404\nPage not found"
		ErrorPage(w, err, 404)
		return
	}
	if r.Method != http.MethodGet {
		err := "405 method not allowed"
		ErrorPage(w, err, 405)
		return
	}

	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}

	slider := r.FormValue("slider")
	radiobtn := r.Form["position"]
	selectvalue := r.FormValue("select")
	range1, _ := strconv.Atoi(r.FormValue("range1"))
	range2, _ := strconv.Atoi(r.FormValue("range2"))
	if range2 < range1 {
		err := "400 Bad Request"
		ErrorPage(w, err, 400)
		return
	}

	res := Filterfunction(slider, radiobtn, selectvalue, range1, range2)
	err = templates.Execute(w, res)
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, 500)
		return
	}
}

func ErrorPage(w http.ResponseWriter, errors string, code int) {
	w.WriteHeader(code)
	ero, err := template.ParseFiles("htmlfiles/error.html")
	if err != nil {
		http.Error(w, "500"+"\n"+"Internal Server Error", 500)
		return
	}
	ero.Execute(w, errors)
}
