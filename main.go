package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/EmiSan1998/gpsTiming-backend/datatypes"
	"github.com/julienschmidt/httprouter"
)

var (
	data map[string]datatypes.Route
)

// ProgramVersion Current version of the microservice
const ProgramVersion string = "0.1.0"

func main() {
	flag.Parse()
	data = make(map[string]datatypes.Route)

	r := httprouter.New()
	r.GET("/route/:key", getRoute)
	r.POST("/route", postRoute)
	r.GET("/statusReport", getStatusReport)

	port := os.Getenv("GPSTIMING_BACKEND_PORT")
	if len(port) == 0 {
		fmt.Println(os.Getenv("GPSTIMING_BACKEND_PORT"))
		port = "8080"
	}

	fmt.Println("gpsTiming-backend")
	fmt.Println("version " + ProgramVersion + " 2020")
	fmt.Println("system online on port " + port)

	addr := flag.String("addr", ":"+port, "http service address")
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func getRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")

	if _, ok := data[k]; ok {
		routeJSON, err := json.Marshal(data[k])
		if err != nil {
			log.Fatal("Serialization error")
		}
		w.Header().Set("Content-Type", "application/json")
		log.Println("Route " + data[k].Name + " loaded " + "from key " + k)
		w.Write(routeJSON)
	} else {
		log.Println("Route not exists")
		w.Write([]byte("Route not exists"))
	}
}

func postRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var route datatypes.Route

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}

	err = json.Unmarshal(jsn, &route)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	randomKey, _ := uuid.NewRandom()
	data[randomKey.String()] = route

	log.Println("Route " + route.Name + " saved " + "on key " + randomKey.String())
	var response struct {
		ID   string
		Name string
	}
	response.ID = randomKey.String()
	response.Name = route.Name

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Write(responseBytes)
}

func getStatusReport(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var response struct {
		BackendVersion string
	}
	response.BackendVersion = ProgramVersion

	log.Println("Status report requested")

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Write(responseBytes)
}
