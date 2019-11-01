package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"os/exec"

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
	saveCaddy()
	restartCaddy()
}

func saveCaddy() {
	tmpl, _ := template.New("Caddy").Parse(`
	{{.Domain}} {
		proxy {{.Source}} {{.Target}}
		tls {{.Email}}
	}`)

	buf := &bytes.Buffer{}
	for _, value := range sites {
		tmpl.Execute(buf, value)
	}

	f, err := os.OpenFile(saveCaddyFile, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}

	//clear file before writing
	f.Truncate(0)
	f.Seek(0, 0)
	// do the actual work
	f.WriteString(buf.String())
}

func startCaddy() {
	cmd := exec.Command("nohup caddy -agree &")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(stdout))
}

func restartCaddy() {
	cmd := exec.Command("pkill -USR1 caddy")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(stdout))
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
	sites[id] = site
	json.NewEncoder(w).Encode(site)
	fmt.Println("Changed Site:")
	fmt.Println(site)
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
	fmt.Println("New Site:")
	fmt.Println(site)
	save()
}

func deleteSite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(params["id"])
	delete(sites, id)
	json.NewEncoder(w).Encode(true)
	fmt.Print("deleted site: ")
	fmt.Println(id)
	save()
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./index.html")
	t.Execute(w, nil)
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

	var dir string

	flag.StringVar(&dir, "dir", "./static", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	router.HandleFunc("/", index).Methods("GET")

	error := http.ListenAndServe(":8000", router)
	save()
	log.Fatal(error)
}

// TODO:
// start caddy with config <- test that
// UI
