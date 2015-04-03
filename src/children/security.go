package children

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Permission represents a security access level for the system. Currently the security access levels are large grain
// but may be fine tuned as development continues
type Permission string

const (
	// Viewer represents a security access level for an individual that is able to only read information from the
	// system, but can not update the information.
	Viewer Permission = "viewer"

	// Modifier represents a security access level for an individual that is able to read and modify information from
	// the system
	Modifier Permission = "modifier"

	// Administrator represents a security access level for an individual that has full access to the system. This
	// includes the ability to read and write information, but also to add new users.
	Administrator Permission = "administrator"
)

var PermissionMap map[Permission]int = map[Permission]int{
	Viewer:        1,
	Modifier:      2,
	Administrator: 3,
}

// Privilege represents a mapping between a specific resource path, such as /child, and the permission on resources in
// that path
type Privilege struct {
	ResourcePath       string     `json:"resourcepath"`
	ResourcePermission Permission `json:"permission"`
}

// Role represents a mapping between a specific individual, Principle, and the privileges to which have been assigned
// to the individual against resources maintained by the system
type Role struct {
	Id         string      `json:"id"`
	Principle  string      `json:"principle"`
	Privileges []Privilege `json:"privileges"`
}

type RoleResource struct{}

func (RoleResource) Id(val interface{}, id string) {
	role := val.(*Role)
	role.Id = id
}

func (RoleResource) Name() string {
	return "Role"
}

func (RoleResource) Path() string {
	return "roles"
}

func (RoleResource) Make(limit int) interface{} {
	found := make([]Role, 0, limit)
	return &found
}

func (RoleResource) New() interface{} {
	return new(Role)
}

func (RoleResource) Validate(val interface{}) bool {
	role := val.(*Role)
	if role.Principle == "" {
		return false
	}
	return true
}

func (RoleResource) UniqueFilter(query *datastore.Query, val interface{}) *datastore.Query {
	role := val.(*Role)
	return query.Filter("Principle=", role.Principle)
}

func (RoleResource) FilterById(query *datastore.Query, id string) *datastore.Query {
	return query
}

func (RoleResource) Index(val interface{}, i int) interface{} {
	return (*val.(*[]Role))[i]
}

func (RoleResource) SingleHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/roles") {
		w.WriteHeader(http.StatusUnauthorized)
	}

	switch r.Method {
	case "GET":
		Fetch(RoleResource{}, c, w, r)
		break
	case "PUT":
		Update(RoleResource{}, c, w, r)
		break
	case "DELETE":
		Delete(RoleResource{}, c, w, r)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		break
	}

}

func (RoleResource) ListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !Authorize(c, w, r, "/roles") {
		w.WriteHeader(http.StatusUnauthorized)
	}

	switch r.Method {
	case "GET":
		Query(RoleResource{}, c, w, r)
		break
	case "POST":
		Create(RoleResource{}, c, w, r)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		break
	}
}

// Authorize validates is the user which initiated the the HTTP request has the required permission level as specifed
// by the point of enforcement (POE) making the authorization request.
func Authorize(c appengine.Context, w http.ResponseWriter, r *http.Request, requiredResourcePath string) bool {

	// Get the current user information from the request. If no user information is present in the request then
	// redirect the request to the login page
	u := user.Current(c)

	if u == nil {

		// TODO - REMOVE - HACK - this is a quick hack to allow me to test and specify the user attempting to access
		// the service via basic authentication. This is a security vulnerability that must be removed before going
		// into production
		if name, _, ok := r.BasicAuth(); ok {
			u = &user.User{Email: name}
		}
	}
	c.Infof("Request by user: %s", u.Email)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return false
	}

	// HACK - hardware a specific user (me) to allways be allowed to access the system with full authority. This is the
	// back door when everything else fails. This really should be removed, but we aren't talking a real high security
	// operation here.
	if u.Email == "davidk.bainbridge@gmail.com" {
		return true
	}

	// Select the required permission based on the operation being performed. We default to Administrator so that the
	// default is the most restrictive permission
	requiredPermission := Administrator
	switch r.Method {
	case "POST":
	case "PUT":
	case "DELETE":
		requiredPermission = Modifier
		break
	case "GET":
		requiredPermission = Viewer
		break
	default:
		requiredPermission = Administrator
		break
	}

	// Looking up the role and privileges of the given user from the data store
	role := make([]Role, 0, 1)
	_, err := datastore.NewQuery("Role").Distinct().Limit(1).Filter("Principle=", u.String()).GetAll(c, &role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// If a role for this individual was found then iterate over the privileges assigned to that principle and check
	// to see if they have been assigned a permission for the required resource path. If they have been assigned a
	// permission make sure that the have the required permission or a permission with more capability.
	if len(role) == 1 {
		for _, privilege := range role[0].Privileges {
			if privilege.ResourcePath == requiredResourcePath {

				has, hasOk := PermissionMap[privilege.ResourcePermission]
				needs, needsOk := PermissionMap[requiredPermission]
				if hasOk && needsOk && has >= needs {
					return true
				}
			}
		}
	}

	// If they don't have the required permission then return false
	return false
}
