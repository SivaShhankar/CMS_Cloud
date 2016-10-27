package Controllers

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	config "github.com/SivaShhankar/CMS_Cloud/Database"
	models "github.com/SivaShhankar/CMS_Cloud/Models"
)

type Users struct {
	C *mgo.Collection
}

func CreateDefaultUserCredentials() {

	AddUsers("sindhuja@ramarson.com", "admin1", true)
	AddUsers("sridharan@ramarson.com", "admin2", true)
	AddUsers("rukmani@ramarson.com", "admin3", true)
	//AddUsers("cloud@ramarson.com", "Admin", true)
}

func ValidateUser(session *mgo.Session, userName, password string) (bool, bool) {

	var tempUsers []models.CredentialsInfo

	dataStore := config.NewDataStore()

	defer dataStore.Close()

	users := retrieveUsers(dataStore)

	iter := users.C.Find(bson.M{"username": userName}).Iter() //, "password": password}).Iter()

	result := models.CredentialsInfo{}
	for iter.Next(&result) {
		tempUsers = append(tempUsers, result)
	}

	if len(tempUsers) > 0 {
		return true, tempUsers[0].AllowToChangePassword
	}

	return false, false
}

func AddUsers(userName, password string, isAllowToChangePassword bool) error {
	dataStore := config.NewDataStore()

	defer dataStore.Close()

	credentials := retrieveUsers(dataStore)

	err := credentials.C.Insert(&models.CredentialsInfo{
		UserName:              userName,
		Password:              password,
		AllowToChangePassword: isAllowToChangePassword,
	})

	return err
}

func retrieveUsers(dataStore *config.DataStore) Users {

	// Gets the current collection
	col := dataStore.Collection("UserCredentials")
	credentials := Users{C: col}

	return credentials
}

func ResetPassword(userName, newPassword string) error {

	dataStore := config.NewDataStore()

	defer dataStore.Close()

	credentials := retrieveUsers(dataStore)

	err := credentials.C.Update(bson.M{"username": userName}, &models.CredentialsInfo{
		UserName:              userName,
		Password:              newPassword,
		AllowToChangePassword: false,
	})

	return err
}
