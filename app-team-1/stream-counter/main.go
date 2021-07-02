package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

const titlepage = `
<html>
<h1>{{ range $i := .}}{{$i.Name}}, {{$i.Location}}, <a href="https://twitter.com/{{$i.Twitter}}">@{{$i.Twitter}}</a><br/>
{{end}}</h1>
</html>`

type Stream struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Twitter  string `json:"twitter"`
}

func handleRequests() {
	log.Println("Starting Streams server...")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/create", create).Methods("POST")
	myRouter.HandleFunc("/list", list)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:5432/postgres", user, password, host))

	defer db.Close()

	if err != nil {
		panic("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}

	db.AutoMigrate(&Stream{})
	handleRequests()
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the stream!")
	fmt.Println("Request Received: /")
}

func create(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var stream Stream
	json.Unmarshal(reqBody, &stream)
	db.Create(&stream)
	fmt.Println("Request Received: /create")
	json.NewEncoder(w).Encode(stream)
}

func list(w http.ResponseWriter, r *http.Request) {
	streams := []Stream{}
	db.Find(&streams)
	fmt.Println("Request Received: /list")
	t := template.Must(template.New("List").Parse(titlepage))
	t.Execute(w, streams)
}
