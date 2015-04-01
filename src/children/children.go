// Package children represents the capabilities to manage information about the children living is group housing at the
// Source of Light home in Port-au-Prince Hatit
package children

import (
	"net/http"

	"appengine"

	"github.com/gorilla/mux"
)

func init() {

	// Set up the HTTP request router for the various resources
	router := mux.NewRouter()

	// Child information
	router.Path("/children/{id}").Methods("GET", "PUT", "DELETE").HandlerFunc(childHandler)
	router.Path("/children").Methods("GET", "POST").HandlerFunc(childrenHandler)

	// Adult information
	router.Path("/adults/{id}").Methods("GET", "PUT", "DELETE").HandlerFunc(adultHandler)
	router.Path("/adults").Methods("GET", "POST").HandlerFunc(adultsHandler)

	// Security information
	router.Path("/security/roles/{id}").Methods("GET", "PUT", "DELETE").HandlerFunc(roleHandler)
	router.Path("/security/roles").Methods("GET", "POST").HandlerFunc(rolesHandler)

	http.Handle("/", router)
}

func roleHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/security/roles") {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func rolesHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/security/roles") {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func adultHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/adults") {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func adultsHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/adults") {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func childHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/children") {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func childrenHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/children") {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
