package main

import (
	"bytes"
	"encoding/json"
	"exec"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

// config

var saveFile = "./sites.json"
var saveCaddyFile = "./Caddyfile"

// models

type Site struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Source      string `json:"source"`
	Target      string `json:"target"`
	Email       string `json:"email"`
}

var sites = make(map[int]Site)

// functions
func getBiggest(map1 map[int]Site) int {
	biggest := 0
	for k := range map1 {
		if k > biggest {
			biggest = k
		}
	}
	return biggest
}

// persist config
func load() {
	x, _ := ioutil.ReadFile(saveFile)
	json.Unmarshal(x, &sites)
}

func save() {
	jsonString, _ := json.Marshal(sites)
	ioutil.WriteFile(saveFile, jsonString, 0777)
}

func saveCaddy() {
	//var caddyFileString = ""
	tmpl, _ := template.New("Caddy").Parse(`
	{{.Domain}} {
		proxy {{.Source}} {{.Target}}
		tls {{.Email}}
	}`)

	buf := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	for _, value := range sites {
		tmpl.Execute(buf, value)
		buf2.WriteString(buf.String())
	}
	fmt.Println("String: ", buf2.String())

	f, err := os.OpenFile(saveCaddyFile, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}

	// do the actual work
	f.WriteString(buf2.String())
}

func startCaddy() {
	cmd := exec.Command("nohup caddy -agree &")
	stdout, err := cmd.Output()
	if err != nil {
		Println(err.Error())
		return
	}

	Print(string(stdout))
}

func restartCaddy() {
	cmd := exec.Command("pkill -USR1 caddy")
	stdout, err := cmd.Output()
	if err != nil {
		Println(err.Error())
		return
	}

	Print(string(stdout))
}

// site controllers
func getSites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

func getSite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(params["id"])
	json.NewEncoder(w).Encode(sites[id])
}

func updateSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	delete(sites, id)
	var site Site
	_ = json.NewDecoder(r.Body).Decode(&site)
	newID := getBiggest(sites) + 1
	site.Id = newID
	sites[newID] = site
	json.NewEncoder(w).Encode(site)
	save()
}

func createSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var site Site
	_ = json.NewDecoder(r.Body).Decode(&site)
	newID := getBiggest(sites) + 1
	site.Id = newID
	sites[newID] = site
	json.NewEncoder(w).Encode(site)
	save()
}

func deleteSite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(params["id"])
	delete(sites, id)
	json.NewEncoder(w).Encode(true)
	save()
}

func main() {
	router := mux.NewRouter()

	load()
	saveCaddy()
	startCaddy()
	//handlers
	router.HandleFunc("/api/Sites", getSites).Methods("GET")
	router.HandleFunc("/api/Sites/{id}", getSite).Methods("GET")
	router.HandleFunc("/api/Sites", createSite).Methods("POST")
	router.HandleFunc("/api/Sites/{id}", updateSite).Methods("PUT")
	router.HandleFunc("/api/Sites/{id}", deleteSite).Methods("DELETE")

	error := http.ListenAndServe(":8000", router)
	save()
	log.Fatal(error)
}

// TODO:
// start caddy with config <- test that
// UI
