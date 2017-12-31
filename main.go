package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"math/rand"
)

type Config struct {
	Port string `yaml:"port"`
}

type CompilationResult struct {
	Output   string   `json:"output"`
	Error    string   `json:"error"`
	Warnings []string `json:"warnings"`
}

type Code struct {
	Code  string `json:"code"`
	Input string `json:"input"`
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
	rand.Seed(42)
	http.HandleFunc("/", homePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		log.Println(string(body))

		code := Code{}
		json.Unmarshal(body, &code)
		log.Println(code.Code)

		response1 := CompilationResult{}
		response1.Output = "lol"
		response1.Warnings = []string{"1: a", "2: b", "3: c"}
		response2 := CompilationResult{}
		response2.Error = "5: Fatality"

		
		i := rand.Intn(2)

		if i%2 == 0 {
			response, _ := json.Marshal(response1)
			fmt.Fprintf(w, string(response))
		} else {
			response, _ := json.Marshal(response2)
			fmt.Fprintf(w, string(response))
		}

	} else {
		t, _ := template.ParseFiles("./templates/index.tpl")
		t.Execute(w, struct{}{})
	}
}
