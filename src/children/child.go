package children

import (
	"net/http"
	"time"

	"appengine"

	"appengine/datastore"
)

// Child represents a child living in group housing
type Child struct {
	Id             string           `json:"id"`
	FamilyName     string           `json:"familyname"`
	GivenName      string           `json:"givenname"`
	BirthDate      time.Time        `json:"birthday,omitempty"`
	EnteredHousing time.Time        `json:"enteredhousing,omitempty"`
	LeftHousing    time.Time        `json:"lefthousing,omitempty"`
	GradeInSchool  int              `json:"gradeinschool,omitempty"`
	Mother         *datastore.Key   `json:"mother,omitempty"`
	Father         *datastore.Key   `json:"father,omitempty"`
	KnownRelatives []*datastore.Key `json:"relatives,omitempty"`
	Siblings       *datastore.Key   `json:"siblings,omitempty"`
}

type ChildResource struct{}

func (ChildResource) Name() string {
	return "Child"
}

func (ChildResource) Id(val interface{}, id string) {
	child := val.(*Child)
	child.Id = id
}

func (ChildResource) Path() string {
	return "children"
}

func (ChildResource) Make(limit int) interface{} {
	found := make([]Child, 0, limit)
	return &found
}

func (ChildResource) New() interface{} {
	return new(Child)
}

func (ChildResource) Validate(val interface{}) bool {
	child := val.(*Child)
	if child.FamilyName == "" || child.GivenName == "" {
		return false
	}
	return true
}

func (ChildResource) UniqueFilter(query *datastore.Query, val interface{}) *datastore.Query {
	child := val.(*Child)
	return query.Filter("FamilyName=", child.FamilyName).Filter("GivenName=", child.GivenName)
}

func (ChildResource) FilterById(query *datastore.Query, id string) *datastore.Query {
	return query.Filter("Id=", id)
}

func (ChildResource) Index(val interface{}, i int) interface{} {
	return (*val.(*[]Child))[i]
}

func (ChildResource) SingleHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/children") {
		w.WriteHeader(http.StatusUnauthorized)
	}

	switch r.Method {
	case "GET":
		Fetch(ChildResource{}, c, w, r)
		break
	case "PUT":
		Update(ChildResource{}, c, w, r)
		break
	case "DELETE":
		Delete(ChildResource{}, c, w, r)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		break
	}
}

func (ChildResource) ListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/children") {
		w.WriteHeader(http.StatusUnauthorized)
	}

	switch r.Method {
	case "GET":
		Query(ChildResource{}, c, w, r)
		break
	case "POST":
		Create(ChildResource{}, c, w, r)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		break
	}
}
