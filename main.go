package main 

import (
	"net/http"
	"log"
	"time"
	"io/ioutil"
	"strings"
)

type Server struct {
	Url		string
	Status	string
}

var serv Server
var ch chan Server

func ping(s * Server) {
	resp, err := http.Get(s.Url)
	if err != nil {
		log.Print(err)
	} else {
		s.Status = resp.Status
	}
	log.Print("Hit server " + s.Status)
}

func start(s *Server) {
	go func () {
		for true {
			ping(s)
			time.Sleep(time.Minute)
		}
	} ()
}

func monitor() {
	serv.Url = "http://www.google.com"
	serv.Status = "not hit"
	ch = make(chan Server)
	start(&serv)
}

func index(res http.ResponseWriter, req *http.Request) {
	log.Print("Status " + serv.Status)
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		res.Write([]byte("Status " + serv.Status))
	} else {
		s_body := "<div class=\"server\">" + serv.Url + "</div><div class=\"server-status\">"+serv.Status+"</div>"
		s_content := strings.Replace(string(content), "<contents>", s_body, 1)
		res.Write([]byte(s_content))		
	}
}

func main() {
	monitor()
	log.Print("Finished starting the monitor process")
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":3000", nil));
}	