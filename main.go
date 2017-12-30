package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"fmt"
)

type Config struct {
	Port    string   `yaml:"port"`
}

var (
	port   = os.Getenv("PORT")
	config Config
)

func init() {
	yamlFile, err := ioutil.ReadFile("./config.yaml")

	if err != nil {
		log.Println("Error while reading config file ", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println("Error while unmarshalling config ", err)
	}

	if port == "" {
		port = config.Port
	}
}

func main() {
	http.HandleFunc("/", homePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	log.Println(string(body))

	if lol := r.PostFormValue("lol"); lol != "" {
		log.Println(lol)
		fmt.Fprint(w, "Hello world!")
	} else {
		t, _ := template.ParseFiles("./templates/index.tpl")
		t.Execute(w, struct{}{})
	}
}
