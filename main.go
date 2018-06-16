package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/phones", getPhonesHandler).Methods("GET")
	router.HandleFunc("/phones/{phone_number}", getPhoneInfoHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getPhonesHandler(response http.ResponseWriter, request *http.Request) {

	keys, ok := request.URL.Query()["userId"]

	if !ok || len(keys) < 1 {
		log.Print("Listando todos los teléfonos")
		return
	}

	userId := keys[0]
	log.Print("Listando los teléfonos del usuario " + userId)
}

func getPhoneInfoHandler(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	log.Print("Listando la información del teléfono " + params["phone_number"])
}
