package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//NewRouter ...
func NewRouter() (*mux.Router, *API) {
	api := API{}

	//Routes ...
	routes := []route{

		route{
			"addPet",
			strings.ToUpper("Post"),
			"/v2/pet/{petId}",
			api.addPet,
		},

		route{
			"deletePet",
			strings.ToUpper("Delete"),
			"/v2/pet/{petId}",
			api.deletePet,
		},

		route{
			"getPets",
			strings.ToUpper("Get"),
			"/v2/pets",
			api.getPets,
		},

		route{
			"getPetByID",
			strings.ToUpper("Get"),
			"/v2/pet/{petId}",
			api.getPetByID,
		},

		route{
			"updatePet",
			strings.ToUpper("Put"),
			"/v2/pet/{petId}",
			api.updatePet,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router, &api
}
