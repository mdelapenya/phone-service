package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/phones", phonesHandler)
	router.HandleFunc("/phones/{phone_number}", phoneByPhoneNumberHandler)

	log.Fatal(http.ListenAndServe(":8000", router))
}

func phonesHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	fmt.Println(vars)
}

func phoneByPhoneNumberHandler(response http.ResponseWriter, request *http.Request) {

}
