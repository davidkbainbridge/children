package children

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/gorilla/mux"
)

type Resource interface {
	Id(val interface{}, id string)
	Name() string
	Path() string
	New() interface{}
	Make(limit int) interface{}
	Index(val interface{}, i int) interface{}
	Validate(val interface{}) bool
	UniqueFilter(query *datastore.Query, val interface{}) *datastore.Query
	FilterById(query *datastore.Query, id string) *datastore.Query

	ListHandler(w http.ResponseWriter, r *http.Request)
	SingleHandler(w http.ResponseWriter, r *http.Request)
}

func Create(resource Resource, c appengine.Context, w http.ResponseWriter, r *http.Request) {

	entityName := resource.Name()
	entityPath := resource.Path()
	child := resource.New()

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &child)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c.Infof("JSON: %s", child)
	if !resource.Validate(child) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := datastore.NewQuery(entityName).Limit(1).Distinct()
	query = resource.UniqueFilter(query, child)
	found := resource.Make(1)
	keys, err := query.GetAll(c, found)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Infof("len = %d", len(keys))
	if len(keys) != 0 {
		http.Error(w, "Duplicate Entry", http.StatusBadRequest)
		return
	}
	parent := datastore.NewKey(c, entityName, entityPath, 0, nil)
	key := datastore.NewIncompleteKey(c, entityName, parent)
	resource.Id(child, uuid.New())
	_, err = datastore.Put(c, key, child)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func Query(resource Resource, c appengine.Context, w http.ResponseWriter, r *http.Request) {
	// Be default the query will return the first 10 items found, but this can be modified by query parameters
	limit := -1
	offset := -1

	entityName := resource.Name()

	// Create new query for the entity type
	query := datastore.NewQuery(entityName).Distinct()

	// The default query string processing of golang does not support order values as well as repeated values. Because
	// of this this code attempt to manually iterate over the query string. This is a fairly primitive parsing of the
	// query string, i.e. split on "&" and then split on "=".
	options := strings.Split(r.URL.RawQuery, "&")
	for _, option := range options {
		nv := strings.SplitN(option, "=", 2)

		// Split option processing into single value (flag) options and name / value options
		if len(nv) == 1 {
			switch nv[0] {
			case "keys":
				query = query.KeysOnly()
				break
			case "count":
				//includeCount = true
				break
			default:
			}
		} else {
			switch nv[0] {
			case "limit":
				i, err := strconv.Atoi(nv[1])
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				query = query.Limit(i)
				limit = i
				break
			case "offset":
				i, err := strconv.Atoi(nv[1])
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				query = query.Offset(i)
				offset = i
				break
			case "filter":
				terms := strings.SplitN(nv[1], ",", 2)
				query = query.Filter(terms[0], terms[1])
				break
			case "order":
				query = query.Order(nv[1])
				break
			default:
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	// If a limit or offset was not specified then set the defaults
	if limit == -1 {
		limit = 100
		query = query.Limit(limit)
	}
	if offset == -1 {
		offset = 0
		query = query.Offset(offset)
	}

	found := resource.Make(limit)
	keys, err := query.GetAll(c, found)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.Infof("FOUND: %d", len(keys))
	encode := json.NewEncoder(w)
	encode.Encode(found)
	w.WriteHeader(http.StatusOK)
}

func Fetch(resource Resource, c appengine.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	entityName := resource.Name()

	query := datastore.NewQuery(entityName).Limit(1).Distinct()
	query = resource.FilterById(query, id)
	found := resource.Make(1)
	keys, err := query.GetAll(c, found)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Infof("len = %d", len(keys))
	if len(keys) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	encode := json.NewEncoder(w)
	encode.Encode(resource.Index(found, 0))
	w.WriteHeader(http.StatusOK)
}

func Update(resource Resource, c appengine.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	entityName := resource.Name()
	child := resource.New()

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &child)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c.Infof("JSON: %s", child)
	if !resource.Validate(child) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := datastore.NewQuery(entityName).Limit(1).Distinct().KeysOnly()
	query = resource.FilterById(query, id)
	keys, err := query.GetAll(c, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(keys) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resource.Id(child, id)
	_, err = datastore.Put(c, keys[0], child)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func Delete(resource Resource, c appengine.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	entityName := resource.Name()

	query := datastore.NewQuery(entityName).Limit(1).Distinct().KeysOnly()
	query = resource.FilterById(query, id)
	keys, err := query.GetAll(c, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(keys) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	datastore.Delete(c, keys[0])
	w.WriteHeader(http.StatusOK)
}
