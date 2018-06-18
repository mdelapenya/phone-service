package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App the application object
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize init the application
func (app *App) Initialize() {
	dbHostname := os.Getenv("DB_HOSTNAME")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connectionString :=
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHostname, user, password, dbName)

	var err error
	app.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := app.DB.Exec(PhonesTableCreationQuery); err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

// Run run the application
func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func (app *App) deletePhoneHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid Phone ID")
		return
	}

	phone := phone{ID: id}
	if err := phone.deletePhone(app.DB); err != nil {
		respondWithError(response, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(response, http.StatusOK, map[string]string{"result": "success"})
}

func (app *App) getPhoneHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid phone ID")
		return
	}

	phone := phone{ID: id}
	if err := phone.getPhone(app.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(response, http.StatusNotFound, "Phone not found")
		default:
			respondWithError(response, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(response, http.StatusOK, phone)
}

func (app *App) getPhonesHandler(response http.ResponseWriter, request *http.Request) {
	count, _ := strconv.Atoi(request.FormValue("count"))
	start, _ := strconv.Atoi(request.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	phones, err := getPhones(app.DB, start, count)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(response, http.StatusOK, phones)
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/phones", app.getPhonesHandler).Methods("GET")
	app.Router.HandleFunc("/phone", app.createPhoneHandler).Methods("POST")
	app.Router.HandleFunc("/phone/{id:[0-9]+}", app.getPhoneHandler).Methods("GET")
	app.Router.HandleFunc("/phone/{id:[0-9]+}", app.updatePhoneHandler).Methods("PUT")
	app.Router.HandleFunc("/phone/{id:[0-9]+}", app.deletePhoneHandler).Methods("DELETE")
}

func main() {
	app := App{}
	app.Initialize()

	app.Run(":8000")
}

func (app *App) createPhoneHandler(response http.ResponseWriter, request *http.Request) {
	var phone phone
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&phone); err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer request.Body.Close()

	if err := phone.createPhone(app.DB); err != nil {
		respondWithError(response, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(response, http.StatusCreated, phone)
}

func respondWithError(response http.ResponseWriter, code int, message string) {
	respondWithJSON(response, code, map[string]string{"error": message})
}

func respondWithJSON(response http.ResponseWriter, code int, payload interface{}) {
	bytes, _ := json.Marshal(payload)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(code)
	response.Write(bytes)
}

func (app *App) updatePhoneHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid phone ID")
		return
	}

	var p phone
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer request.Body.Close()
	p.ID = id

	if err := p.updatePhone(app.DB); err != nil {
		respondWithError(response, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(response, http.StatusOK, p)
}
