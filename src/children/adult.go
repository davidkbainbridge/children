package children

import (
	"net/http"
	"time"

	"appengine"

	"appengine/datastore"
)

// ContactInfo defines basic contact information for an individual. This information is purposfully left vague as
// I am not quite sure what we need to capture for addresses in Haiti.
type ContactInfo struct {
	Line1  string `json:"line1.omitempty"`
	Line2  string `json:"line2.omitempty"`
	Line3  string `json:"line3.omitempty"`
	Line4  string `json:"line4.omitempty"`
	Phone1 string `json:"phone1.omitempty"`
	Phone2 string `json:"phone2.omitempty"`
}

// Adult represents a responsible person of the age of majority and is used to help track any known relatives of one
// of the children in group housing
type Adult struct {
	Id               string      `json:"id"`
	FamilyName       string      `json:"familyname"`
	GivenName        string      `json:"givenname"`
	Relationship     string      `json:"relationship,omitempty"`
	LastContact      time.Time   `json:"lastcontact,omitempty"`
	LastVisit        time.Time   `json:"lastvisit,omitempty"`
	LastKnownContact ContactInfo `json:"lastknowncontact,omitempty"`
}

type AdultResource struct{}

func (AdultResource) Id(val interface{}, id string) {
	adult := val.(*Adult)
	adult.Id = id
}

func (AdultResource) Name() string {
	return "Adult"
}

func (AdultResource) Path() string {
	return "adults"
}

func (AdultResource) Make(limit int) interface{} {
	found := make([]Adult, 0, limit)
	return &found
}

func (AdultResource) New() interface{} {
	return new(Adult)
}

func (AdultResource) Validate(val interface{}) bool {
	Adult := val.(*Adult)
	if Adult.FamilyName == "" || Adult.GivenName == "" {
		return false
	}
	return true
}

func (AdultResource) UniqueFilter(query *datastore.Query, val interface{}) *datastore.Query {
	Adult := val.(*Adult)
	return query.Filter("FamilyName=", Adult.FamilyName).Filter("GivenName=", Adult.GivenName)
}

func (AdultResource) FilterById(query *datastore.Query, id string) *datastore.Query {
	return query
}

func (AdultResource) Index(val interface{}, i int) interface{} {
	return (*val.(*[]Adult))[i]
}

func (AdultResource) SingleHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/adults") {
		w.WriteHeader(http.StatusUnauthorized)
	}

	switch r.Method {
	case "GET":
		Fetch(AdultResource{}, c, w, r)
		break
	case "PUT":
		Update(AdultResource{}, c, w, r)
		break
	case "DELETE":
		Delete(AdultResource{}, c, w, r)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		break
	}
}

func (AdultResource) ListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/adults") {
		w.WriteHeader(http.StatusUnauthorized)
	}

	switch r.Method {
	case "GET":
		Query(AdultResource{}, c, w, r)
		break
	case "POST":
		Create(AdultResource{}, c, w, r)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		break
	}
}
