package Database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
)

type configuration struct {
	Server,
	MongoDBHost,
	DBUser,
	DBPwd,
	Database,
	CVLocation string
}

var AppConfig configuration

func LoadAppConfig() {
	file, err := os.Open("Database/config.json")
	defer file.Close()
	if err != nil {
		log.Fatalf("[LoadAppConfig]: %s\n", err)
	}
	decoder := json.NewDecoder(file)
	AppConfig = configuration{}
	err = decoder.Decode(&AppConfig)

	if err != nil {
		log.Fatalf("[LoadAppConfig]: %s\n", err)
	}
}

var Session *mgo.Session

func CreateDBSession() {
	var err error
	// Session, err = mgo.DialWithInfo(&mgo.DialInfo{
	// 	Addrs:    []string{AppConfig.MongoDBHost},
	// 	Username: AppConfig.DBUser,
	// 	Password: AppConfig.DBPwd,
	// 	Timeout:  60 * time.Second,
	// })
	Session, err = mgo.Dial("mongodb://gowri:gowri@ds035796.mlab.com:35796/mycmstool")
	if err != nil {
		fmt.Println(err)
		log.Fatalf("[CreateDbSession]: %s\n", err)
	}
}

func AddIndexes() {
	var err error
	userIndex := mgo.Index{
		Key:        []string{"mobile", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	// Add indexes into MongoDB
	session := getSession().Copy()
	defer session.Close()
	userCol := session.DB(AppConfig.Database).C("JobCandidates")

	err = userCol.EnsureIndex(userIndex)
	if err != nil {
		log.Fatalf("[AddIndexes]: %s\n", err)
	}
}

func AddCredentialIndexes() {

	var err error
	userIndex := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	// Add indexes into MongoDB
	session := getSession().Copy()
	defer session.Close()
	userCol := session.DB(AppConfig.Database).C("UserCredentials")

	err = userCol.EnsureIndex(userIndex)
	if err != nil {
		log.Fatalf("[AddIndexes]: %s\n", err)
	}
}

func getSession() *mgo.Session {
	if Session == nil {
		var err error
		// Session, err = mgo.DialWithInfo(&mgo.DialInfo{
		// 	Addrs:    []string{AppConfig.MongoDBHost},
		// 	Username: AppConfig.DBUser,
		// 	Password: AppConfig.DBPwd,
		// 	Timeout:  60 * time.Second,
		// })
		Session, err = mgo.Dial("mongodb://gowri:gowri@ds035796.mlab.com:35796/mycmstool")
		if err != nil {
			log.Fatalf("[getSession]: %s\n", err)
		}
	}
	return Session
}

// DataStore for MongoDB
type DataStore struct {
	MongoSession *mgo.Session
}

// Close closes a mgo.Session value.
func (ds *DataStore) Close() {
	ds.MongoSession.Close()
}

// Collection returns mgo.collection for the given name
func (ds *DataStore) Collection(name string) *mgo.Collection {
	return ds.MongoSession.DB(AppConfig.Database).C(name)
}

// NewDataStore creates a new DataStore object to be used for each HTTP request.
func NewDataStore() *DataStore {
	session := getSession().Copy()
	dataStore := &DataStore{
		MongoSession: session,
	}
	return dataStore
}
