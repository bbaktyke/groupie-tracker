package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Index struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Loc struct {
	Ind []Index `json:"index"`
}

type RelationStruct struct {
	DatesLocations map[string][]string
}

type Everything struct {
	Everyone []Artist
	Location Loc
}

type General struct {
	Artist

	RelationStruct
}

func GetAllLocations() Loc {
	var location Loc
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return location
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return location
	}
	err = json.Unmarshal(bytes, &location)
	if err != nil {
		return location
	}
	defer response.Body.Close()
	return location
}

func MainParser() []Artist {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var artists []Artist

	json.Unmarshal(body, &artists)

	return artists
}

func RelationParser(id int) RelationStruct {
	res, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", id))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var r RelationStruct

	json.Unmarshal(body, &r)

	return r
}

func idParser(id int) Artist {
	res, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%d", id))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var a Artist

	json.Unmarshal(body, &a)

	return a
}

func Finder(s string) []Artist {
	splitted := strings.Split(s, " - ")
	trimmed := strings.TrimSpace(splitted[0])
	tittled := strings.Title(trimmed)
	capsed := strings.ToUpper(trimmed)
	artists := MainParser()
	location := GetAllLocations()

	m := make(map[int]int)

	for i := 0; i < len(artists); i++ {
		if strings.Contains(artists[i].Name, trimmed) || strings.Contains(artists[i].Name, tittled) || strings.Contains(artists[i].Name, capsed) {
			m[artists[i].ID]++
		}

		for j := 0; j < len(artists[i].Members); j++ {
			if strings.Contains(artists[i].Members[j], trimmed) || strings.Contains(artists[i].Members[j], tittled) || strings.Contains(artists[i].Members[j], capsed) {
				m[artists[i].ID]++
			}
		}

		if strings.Contains(strconv.Itoa(artists[i].CreationDate), trimmed) {
			m[artists[i].ID]++
		}

		if strings.Contains(artists[i].FirstAlbum, trimmed) {
			m[artists[i].ID]++
		}

	}

	for _, val := range location.Ind {
		for _, v := range val.Locations {
			if strings.Contains(strings.ToLower(v), trimmed) {
				m[val.ID]++
			}
		}
	}

	var res []Artist

	for key, _ := range m {
		for _, val := range artists {
			if val.ID == key {
				res = append(res, val)
				break
			}
		}
	}
	return res
	// fmt.Println(m)
}

func Filterfunction(s string, x []string, selectval string, range1 int, range2 int) []Artist {
	converted, err := strconv.Atoi(s)
	location := GetAllLocations()
	if err != nil {
		fmt.Println("Error not valid data")
	}
	var sliceofchx []int
	if len(x) > 0 {
		for _, v := range x {
			y, _ := strconv.Atoi(v)
			sliceofchx = append(sliceofchx, y)
		}
	}

	var m []int
	artists := MainParser()
	for i := 0; i < len(artists); i++ {
		splitalbum := strings.Split(artists[i].FirstAlbum, "-")
		convfirstalbum, _ := strconv.Atoi(splitalbum[2])
		if len(x) > 0 {
			for _, v := range sliceofchx {
				if artists[i].CreationDate >= converted && v == len(artists[i].Members) && (convfirstalbum >= range1 && convfirstalbum <= range2) {
					m = append(m, artists[i].ID)
				}
			}
		} else {
			if artists[i].CreationDate >= converted && convfirstalbum >= range1 && convfirstalbum <= range2 {
				m = append(m, artists[i].ID)
			}
		}
	}

	var r []int
	var inter []int
	var res []Artist
	if selectval != "" {
		for _, val := range location.Ind {
			for _, v := range val.Locations {
				if v == selectval {
					r = append(r, val.ID)
				}
			}
		}
	} else {
		r = m
	}
	inter = Inter(m, r)

	for _, key := range inter {
		for _, val := range artists {
			if val.ID == key {
				res = append(res, val)
				break
			}
		}
	}
	return res
}

func Inter(w1, w2 []int) []int {
	var x []int
	for i := 0; i < len(w1); i++ {
		for j := 0; j < len(w2); j++ {
			if w1[i] == w2[j] {
				x = append(x, w1[i])
			}
		}
	}
	return x
}
