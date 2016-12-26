package Models

import "gopkg.in/mgo.v2/bson"

type (
	ApplicantInfo struct {
		ID             bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Name           string        `json:"name"`
		DateOfBirth    string     	 `json:"dateOfBirth"`
		Age            int           `json:"age"`
		Gender         string        `json:"gender"`
		Mobile         int           `json:"mobile"`
		AlternativeMobile int		 `json:"alternativeMobile"`
		Email          string        `json:"email"`
		Location       string        `json:"location"`
		Qualification  string        `json:"qualification"`
		Specialization string        `json:"specialization"`
		Department     string        `json:"department"`
		JobCode        string        `json:"jobcode"`
		JobApplicationDate string `json:"jobApplicationDate"`
		Position       string        `json:"position"`
		Experience     float64       `json:"experience"`
		CvPath         string        `json:"cvpath"`
		SourceFrom     string        `json:"sourcefrom"`
		Impression	   string		 `json:"impression"`
		Result		   string		 `json:"result"`
		CloudObject    string        `json:"cloudobject"`
	}

	// UserInfo -Capture logged user information
	UserInfo struct {
		ID           string `json:"id"`
		EMail        string `json:"email"`
		VerifiedMail bool   `json:"verified_email"`
		Name         string `json:"name"`
		GivenName    string `json:"given_name"`
		FamilyName   string `json:"family_name"`
		Pciture      string `json:"pciture"`
		Locale       string `json:"locale"`
		HD           string `'json:"hd"`
	}
)
