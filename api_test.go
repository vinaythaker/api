package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var a App
var jwtToken *http.Cookie

func TestMain(t *testing.T) {
	a = App{}
	a.Initialize()
	a.Run()
}

func TestMalformedRequestBodyAddPet(t *testing.T) {

	payload := []byte(``)
	req, _ := http.NewRequest("POST", "/v2/pet/10", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request payload" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
	}
}

func TestMalformedRequestURLAddPet(t *testing.T) {

	payload := []byte(`{"id":10,"name":"grrr"}`)
	req, _ := http.NewRequest("POST", "/v2/pet/8&6%434", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid pet ID'. Got '%s'", m["error"])
	}
}

func TestAddPet(t *testing.T) {

	payload := []byte(`{"id":1,"name":"grrr"}`)
	req, _ := http.NewRequest("POST", "/v2/pet/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var pet Pet
	json.Unmarshal(response.Body.Bytes(), &pet)

	if pet.Name != "grrr" {
		t.Errorf("Expected user name to be 'grrr'. Got '%v'", pet.Name)
	}
}

func TestMalformedRequestURLGetPet(t *testing.T) {

	req, _ := http.NewRequest("GET", "/v2/pet/8&6%434", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid pet ID'. Got '%s'", m["error"])
	}
}

func TestGetNonExistentPet(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v2/pet/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Pet not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Pet not found'. Got '%s'", m["error"])
	}
}

func TestGetPet(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v2/pet/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetPets(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v2/pets", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var pets []Pet
	json.Unmarshal(response.Body.Bytes(), &pets)
	if len(pets) != 1 {
		t.Errorf("Expected 1 pet got %d", len(pets))
	}
}

func TestMalformedRequestBodyUpdatePet(t *testing.T) {

	payload := []byte(``)
	req, _ := http.NewRequest("PUT", "/v2/pet/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request payload" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
	}
}

func TestMalformedRequestURLUpdatePet(t *testing.T) {

	payload := []byte(`{"id":1,"name":"smile"}`)
	req, _ := http.NewRequest("PUT", "/v2/pet/%21%40#23%24%25%5E1234567890abcdefghijklmnopqrstuvwxyz", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid pet ID'. Got '%s'", m["error"])
	}
}

func TestUpdateNonExistentPet(t *testing.T) {

	payload := []byte(`{"id":45,"name":"smile"}`)

	req, _ := http.NewRequest("PUT", "/v2/pet/45", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Pet not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Pet not found'. Got '%s'", m["error"])
	}
}

func TestUpdatePet(t *testing.T) {

	payload := []byte(`{"id":1,"name":"grrr"}`)
	req, _ := http.NewRequest("POST", "/v2/pet/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	payload = []byte(`{"id":1,"name":"smile"}`)
	req, _ = http.NewRequest("PUT", "/v2/pet/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var pet Pet
	json.Unmarshal(response.Body.Bytes(), &pet)

	if pet.Name != "smile" {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", "grrr", "smile", pet.Name)
	}
}

func TestMalformedRequestURLDeletePet(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/v2/pet/%21%40#23%24%25%5E1234567890abcdefghijklmnopqrstuvwxyz", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid pet ID'. Got '%s'", m["error"])
	}
}

func TestDeleteNonExistentPet(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/v2/pet/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Pet not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Pet not found'. Got '%s'", m["error"])
	}
}

func TestDeletePet(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/v2/pet/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/v1/pet/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestShutdown(t *testing.T) {

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	a.server.Shutdown(ctx)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	if jwtToken != nil {
		req.AddCookie(jwtToken)
	}
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
