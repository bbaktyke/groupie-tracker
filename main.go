package main

import (
	"bbaktyke/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("http://localhost:8000")
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", server.MainPage)
	http.HandleFunc("/artists/", server.ArtistPage)
	http.HandleFunc("/search/", server.SearchHadler)
	http.HandleFunc("/filter/", server.FilterHadler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
