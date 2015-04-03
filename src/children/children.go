// Package children represents the capabilities to manage information about the children living is group housing at the
// Source of Light home in Port-au-Prince Hatit
package children

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {

	// Set up the HTTP request router for the various resources
	router := mux.NewRouter()

	// Child information
	router.Path("/children/{id}").Methods("GET", "PUT", "DELETE").HandlerFunc(ChildResource{}.SingleHandler)
	router.Path("/children").Methods("GET", "POST").HandlerFunc(ChildResource{}.ListHandler)

	// Adult information
	router.Path("/adults/{id}").Methods("GET", "PUT", "DELETE").HandlerFunc(AdultResource{}.SingleHandler)
	router.Path("/adults").Methods("GET", "POST").HandlerFunc(AdultResource{}.ListHandler)

	// Security information
	router.Path("/security/roles/{id}").Methods("GET", "PUT", "DELETE").HandlerFunc(RoleResource{}.SingleHandler)
	router.Path("/security/roles").Methods("GET", "POST").HandlerFunc(RoleResource{}.ListHandler)

	http.Handle("/", router)
}
