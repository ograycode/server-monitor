package main 

import (
	"net/http"
	"log"
	"time"
	"io/ioutil"
	"strings"
)

/*	Server Struct
*	Basic server data with Url and Status
*/
type Server struct {
	Url		string
	Status	string
}

/* Server map type which holds the server struct */
type serverMap map[string] Server

/* Global package level server map */
var sMap serverMap

/* Ping given servers in the server map*/
func ping() {
	for _, v := range sMap {
		go func () {
			log.Print("Hitting..." + v.Url)
			resp, err := http.Get(v.Url)
			if err != nil {
				log.Print(err)
			} else {
				v.Url = resp.Status
			}
			log.Print("Hit server " + v.Status)
		} ()
	}
}

/* Start the background thread to ping server */
func start() {
	go func () {
		for true {
			log.Print("Starting Ping")
			ping()
			time.Sleep(time.Minute)
		}
	} ()
}

/* Set up the monitor thread */
func monitor() {
	serv := new(Server)
	serv.Url = "http://www.google.com"
	serv.Status = "not hit"
	sMap = make(serverMap)
	sMap[serv.Url] = *serv
	start()
}

/*Handles the response to "/" */
func index(res http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal("index.html not found")
	} else {
		s_body := ""
		for _,v := range sMap {
			s_body += "<div class=\"server\">" + v.Url + "</div>"
			s_body += "<div class=\"server-status\">"+v.Status+"</div>"
		}
		
		s_content := strings.Replace(string(content), "<contents>", s_body, 1)
		res.Write([]byte(s_content))		
	}
}

/*Start of the program */
func main() {
	monitor()
	log.Print("Finished starting the monitor process")
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":3000", nil));
}	