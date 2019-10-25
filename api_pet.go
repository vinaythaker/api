package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API ...
type API struct {
	petMap map[int64]Pet
	app    *App
}

func (a *API) getPets(w http.ResponseWriter, r *http.Request) {

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	p := Pet{}
	pets, err := p.getPets(a.app.db, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, pets)

}

func (a *API) getPetByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["petId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	p := Pet{ID: int64(id)}
	err = p.getPet(a.app.db)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Pet not found")
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *API) addPet(w http.ResponseWriter, r *http.Request) {
	var p Pet

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["petId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := p.addPet(a.app.db); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		respondWithJSON(w, http.StatusCreated, p)
		return
	}
}

func (a *API) deletePet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["petId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	p := Pet{ID: int64(id)}
	if err := p.deletePet(a.app.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Pet not found")
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}
}

func (a *API) updatePet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["petId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var p Pet
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.updatePet(a.app.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Pet not found")
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		respondWithJSON(w, http.StatusOK, p)
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}
