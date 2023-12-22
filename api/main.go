package main

import (
	"log"
	"net/http"
	"projet/api/REST"
	"strconv"
)

const PORT = 443

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"

func Handler(w http.ResponseWriter, r *http.Request) {

	var c REST.ApiController
	c.Dispatch(r, w)
}

func main() {

	http.HandleFunc("/", Handler)

	log.Println(Green, "Starting local server at port", PORT, Reset)

	if err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil); err != nil {
		log.Fatal(Red, err, Reset)
	}
}
