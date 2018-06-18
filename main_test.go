package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var app App

func TestCreatePhone(t *testing.T) {
	clearTable()

	payload := []byte(`{"phone":"666111222","company":"Movistar","phoneType="mobile","userId":"1"}`)

	req, _ := http.NewRequest("POST", "/phone", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["phone"] != "666111222" {
		t.Errorf("Expected phone number to be '666111222'. Got '%v'", m["phone"])
	}

	if m["company"] != "Movistar" {
		t.Errorf("Expected phone company to be 'Movistar'. Got '%v'", m["company"])
	}

	if m["phoneType"] != "mobile" {
		t.Errorf("Expected phone type to be 'mobile'. Got '%v'", m["mobile"])
	}

	if m["userId"] != "1" {
		t.Errorf("Expected userId to be '1'. Got '%v'", m["userId"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected ID to be '1'. Got '%v'", m["id"])
	}
}

func TestDeletePhone(t *testing.T) {
	clearTable()
	addPhones(1)

	req, _ := http.NewRequest("GET", "/phone/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/phone/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/phone/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/phones", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetPhone(t *testing.T) {
	clearTable()
	addPhones(1)

	req, _ := http.NewRequest("GET", "/phone/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetNonExistentPhone(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/phone/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string

	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Phone not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Phone not found'. Got '%s'", m["error"])
	}
}

func TestMain(m *testing.M) {
	app = App{}

	app.Initialize()

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func TestUpdatePhone(t *testing.T) {
	clearTable()
	addPhones(1)

	req, _ := http.NewRequest("GET", "/phone/1", nil)
	response := executeRequest(req)
	var originalPhone map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalPhone)

	payload := []byte(`{"phone":"666111222","company":"Movistar","phoneType="updated","userId":"1"}`)

	req, _ = http.NewRequest("PUT", "/phone/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["phone"] != originalPhone["phone"] {
		t.Errorf("Expected the phone to remain the same (%v). Got %v",
			originalPhone["phone"], m["phone"])
	}

	if m["phoneType"] == originalPhone["phoneType"] {
		t.Errorf(
			"Expected the phone type to change from '%v' to '%v'. Got '%v'",
			originalPhone["phoneType"], "updated", m["phoneType"])
	}
}

func addPhones(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		index := strconv.Itoa(i)

		app.DB.Exec(
			"INSERT INTO phones(phone, company, phoneType, userId) VALUES($1, $2, $3, $4)",
			"Phone "+index, "Company"+index, "mobile", string(i))
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM phones")
	app.DB.Exec("ALTER SEQUENCE phones_id_seq RESTART WITH 1")
}

func ensureTableExists() {
	if _, err := app.DB.Exec(PhonesTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, request)

	return recorder
}
