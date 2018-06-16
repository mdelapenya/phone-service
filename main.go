package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type phoneResource struct {
	Phone     string `json:"phone,omitempty"`
	Company   string `json:"company,omitempty"`
	PhoneType string `json:"phone_type,omitempty"`
	UserID    string `json:"userId,omitempty"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/phones", getPhonesHandler).Methods("GET")
	router.HandleFunc("/phones/{phone}", getPhoneInfoHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getPhonesHandler(response http.ResponseWriter, request *http.Request) {
	var phones []phoneResource

	phones = append(phones, phoneResource{Phone: "123456789", Company: "Movistar", PhoneType: "Fijo", UserID: "242"})
	phones = append(phones, phoneResource{Phone: "987654321", Company: "Orange", PhoneType: "Móvil", UserID: "1234"})
	phones = append(phones, phoneResource{Phone: "234293735", Company: "Vodafone", PhoneType: "Móvil", UserID: "3242"})

	keys, ok := request.URL.Query()["userId"]

	if !ok || len(keys) < 1 {
		log.Print("Listando todos los teléfonos")
		json.NewEncoder(response).Encode(phones)
		return
	}

	userID := keys[0]
	log.Print("Listando los teléfonos del usuario " + userID)
}

func getPhoneInfoHandler(response http.ResponseWriter, request *http.Request) {
	var phones []phoneResource

	phones = append(phones, phoneResource{Phone: "123456789", Company: "Movistar", PhoneType: "Fijo", UserID: "242"})
	phones = append(phones, phoneResource{Phone: "987654321", Company: "Orange", PhoneType: "Móvil", UserID: "1234"})
	phones = append(phones, phoneResource{Phone: "234293735", Company: "Vodafone", PhoneType: "Móvil", UserID: "3242"})

	params := mux.Vars(request)
	for _, item := range phones {
		if item.Phone == params["phone"] {
			json.NewEncoder(response).Encode(item)
			return
		}
	}
	response.WriteHeader(http.StatusNotFound)
	json.NewEncoder(response).Encode("Phone not found")
	log.Print("Listando la información del teléfono " + params["phone_number"])
}
