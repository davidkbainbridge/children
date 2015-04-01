package children

import (
	"time"

	"appengine/datastore"
)

// ContactInfo defines basic contact information for an individual. This information is purposfully left vague as
// I am not quite sure what we need to capture for addresses in Haiti.
type ContactInfo struct {
	Line1  string
	Line2  string
	Line3  string
	Line4  string
	Phone1 string
	Phone2 string
}

// Adult represents a responsible person of the age of majority and is used to help track any known relatives of one
// of the children in group housing
type Adult struct {
	FamilyName       string
	GivenName        string
	Relationship     string
	LastContact      time.Time
	LastVisit        time.Time
	LastKnownContact ContactInfo
}

// Child represents a child living in group housing
type Child struct {
	FamilyName     string
	GivenName      string
	BirthDate      time.Time
	EnteredHousing time.Time
	LeftHousing    time.Time
	GradeInSchool  int
	Mother         *datastore.Key
	Father         *datastore.Key
	KnownRelatives []*datastore.Key
	Siblings       *datastore.Key
}
