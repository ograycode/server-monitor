package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

/*	Server Struct
*	Basic server data with Url and Status
 */
type Server struct {
	Url    string
	Status string
}

/* Server map type which holds the server struct */
type serverMap map[string]Server

/* Global package level server map */
var sMap serverMap

/* Ping given servers in the server map*/
func ping() {
	for k, v := range sMap {
		go func(k string, v Server) {
			log.Print("Hitting..." + v.Url)
			resp, err := http.Get(v.Url)
			if err != nil {
				log.Print(err)
			} else {
				v.Status = resp.Status
				sMap[k] = v
			}
			log.Print("Hit server " + v.Status)
		}(k, v)
	}
}

/* Start the background thread to ping server */
func start() {
	go func() {
		for true {
			log.Print("Starting Ping")
			ping()
			time.Sleep(time.Minute)
		}
	}()
}

/* Set up the monitor thread */
func monitor() {
	sMap = make(serverMap)
	var servers []Server
	content, err := ioutil.ReadFile("servers.json")
	if err != nil {
		log.Fatal("Unable to find servers.json ", err)
	}
	err = json.Unmarshal(content, &servers)
	if err != nil {
		log.Fatal("JSON not formatted correctly, real err: ", err)
	}
	log.Print(servers)
	for i := range servers {
		sMap[servers[i].Url] = servers[i]
	}
	log.Print(sMap)
	start()
}

/*Handles all server requests and routes to the correct controller */
func router(res http.ResponseWriter, req *http.Request) {

	log.Print("Method: " + req.Method)
	log.Print("Url:" + req.URL.Path)

	switch {
	case req.URL.Path == "/", req.URL.Path == "/index":
		processIndexRequest(res, req)
	case req.URL.Path == "/favicon.ico":
		processFaviconRequest(res, req)
	case strings.Contains(req.URL.Path, "/public"):
		http.Redirect(res, req, "/", http.StatusFound)
	default:
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

//Processes the index request as forwarded by the router
func processIndexRequest(res http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal("index.html not found")
	} else {
		s_body := ""
		for _, v := range sMap {
			s_body += "<div>"
			s_body += "<span class=\"server\">" + v.Url + "</span>"
			s_body += "<span class=\"server-status\">" + v.Status + "</span>"
			s_body += "</div>"
		}
		res.Write([]byte(strings.Replace(string(content), "<contents>", s_body, 1)))
	}
}

//Proceses the favicon.ico request as forwarded by the controller
func processFaviconRequest(res http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile("favicon.ico")
	if err != nil {
		log.Fatal("favicon.ico not found ", err)
	}
	res.Write(content)
}

/*Start of the program */
func main() {
	monitor()
	log.Print("Finished starting the monitor process")
	http.HandleFunc("/", router)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
