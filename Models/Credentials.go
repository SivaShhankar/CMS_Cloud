package Models

import "gopkg.in/mgo.v2/bson"

type CredentialsInfo struct {
	ID                    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserName              string        `json:"username"`
	Password              string        `json:"password"`
	AllowToChangePassword bool          `json:"allowtochangepassword"`
}
