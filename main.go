package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
)

type Pet struct {
	Name string
	Age  int
}

var Pets []Pet

const servePort int = 8080

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func main() {
	http.HandleFunc("/pet", petHandler)
	http.HandleFunc("/health", healthCheckHandler)

	log.Println("server start")
	log.Println("server IP:", getLocalIP())

	address := ":" + strconv.Itoa(servePort)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func petHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPet(w, r)
	case http.MethodPost:
		postPet(w, r)
	default:
		http.Error(w, "http уйню сморозил", http.StatusMethodNotAllowed)
	}
}

func getPet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Pets)

	log.Println("client:", r.RemoteAddr)
}

func postPet(w http.ResponseWriter, r *http.Request) {
	var pet Pet
	err := json.NewDecoder(r.Body).Decode(&pet)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Pets = append(Pets, pet)

	fmt.Fprintf(w, "post new pet: '%v'", pet)

	log.Println("client:", r.RemoteAddr)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "всё заебись...")
}
