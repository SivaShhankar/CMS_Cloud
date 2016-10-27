package Controllers

import (
	"fmt"
	"net/http"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"

	cloudStorage "github.com/SivaShhankar/CMS_Cloud/CloudStorage"
	config "github.com/SivaShhankar/CMS_Cloud/Database"
	models "github.com/SivaShhankar/CMS_Cloud/Models"
)

func retrieveCurrentCollection(dataStore *config.DataStore) Candidates {

	// Gets the current collection
	col := dataStore.Collection("JobCandidates")
	candidates := Candidates{C: col}

	return candidates
}

func GetAllApplicantsInfo(session *mgo.Session) []models.ApplicantInfo {

	var applicants []models.ApplicantInfo

	// create new data store.
	dataStore := config.NewDataStore()

	// Close the session.
	defer dataStore.Close()

	candidates := retrieveCurrentCollection(dataStore)

	iter := candidates.C.Find(nil).Sort("name").Iter()

	result := models.ApplicantInfo{}

	for iter.Next(&result) {
		applicants = append(applicants, result)
	}

	return applicants
}

func GetApplicantByMobileNumber(session *mgo.Session, MobileNumber int) []models.ApplicantInfo {

	var applicants []models.ApplicantInfo

	dataStore := config.NewDataStore()

	defer dataStore.Close()

	candidates := retrieveCurrentCollection(dataStore)

	iter := candidates.C.Find(bson.M{"mobile": MobileNumber}).Iter()

	result := models.ApplicantInfo{}

	for iter.Next(&result) {

		applicants = append(applicants, result)
	}

	return applicants
}

type Candidates struct {
	C *mgo.Collection
}

func SearchCandidatesByType(session *mgo.Session, searchType, searchValue string) []models.ApplicantInfo {

	var applicants []models.ApplicantInfo

	dataStore := config.NewDataStore()

	defer dataStore.Close()

	candidates := retrieveCurrentCollection(dataStore)

	iter := candidates.C.Find(bson.M{searchType: &bson.RegEx{Pattern: searchValue, Options: "i"}}).Sort("name").Iter()

	result := models.ApplicantInfo{}

	for iter.Next(&result) {
		applicants = append(applicants, result)
	}

	return applicants
}

func FilterCandidatesByRange(session *mgo.Session, filterType, rangeFrom, rangeTo string) []models.ApplicantInfo {

	var applicants []models.ApplicantInfo

	from, _ := strconv.Atoi(rangeFrom)
	to, _ := strconv.Atoi(rangeTo)

	dataStore := config.NewDataStore()

	defer dataStore.Close()

	candidates := retrieveCurrentCollection(dataStore)

	iter := candidates.C.Find(bson.M{filterType: bson.M{"$gte": from, "$lte": to}}).Sort(filterType).Iter()

	result := models.ApplicantInfo{}

	for iter.Next(&result) {
		applicants = append(applicants, result)
	}

	return applicants
}

func DeleteCandidateByMobileNumber(session *mgo.Session, mobileNumber string) error {

	mobile, _ := strconv.Atoi(mobileNumber)

	fmt.Println("Mobile Number: ", mobile)

	dataStore := config.NewDataStore()

	defer dataStore.Close()

	candidates := retrieveCurrentCollection(dataStore)

	err := candidates.C.Remove(bson.M{"mobile": mobile})

	return err
}

func SaveInfo(session *mgo.Session, r *http.Request, mode string) {

	var err error
	name := r.FormValue("name")
	sage := r.FormValue("age")
	gender := r.FormValue("gender")

	sOldMobile := r.FormValue("oldMobile")
	smobile := r.FormValue("mobile")
	email := r.FormValue("email")
	location := r.FormValue("location")

	qualification := r.FormValue("qualification")
	specialization := r.FormValue("specialization")
	department := r.FormValue("department")
	jobCode := r.FormValue("jobCode")
	position := r.FormValue("position")
	sExpMonth := r.FormValue("expMonth")
	sExpYear := r.FormValue("expYear")
	sourceFrom := r.FormValue("sourceFrom")

	age, _ := strconv.Atoi(sage)
	mobile, _ := strconv.Atoi(smobile)
	OldMobile, _ := strconv.Atoi(sOldMobile)
	sExperience := sExpYear + "." + sExpMonth

	experience, _ := strconv.ParseFloat(sExperience, 64)

	fmt.Println(experience)

	_, handler, err := r.FormFile("file")

	var StoragePath, CloudObject string

	// If no file has selected in the Form, it will throw an error
	// Cond 1 : if mode  is update, then retreive file value from hidden text box
	// Cond 2 : if mode is Insert, then retreive file value from file field

	if err != nil && mode == "Update" {
		StoragePath = r.FormValue("uploadedFile")
		CloudObject = r.FormValue("cloudobject")
	} else {
		StoragePath, CloudObject = cloudStorage.GCloudUploadFiles(r, r.FormValue("name")+"-"+r.FormValue("mobile")+"-"+handler.Filename)
	}
	dataStore := config.NewDataStore()

	defer dataStore.Close()

	candidates := retrieveCurrentCollection(dataStore)

	if mode == "Insert" {
		err = candidates.C.Insert(&models.ApplicantInfo{
			Name:           name,
			Age:            age,
			Gender:         gender,
			Mobile:         mobile,
			Email:          email,
			Location:       location,
			Qualification:  qualification,
			Specialization: specialization,
			Department:     department,
			JobCode:        jobCode,
			Position:       position,
			Experience:     experience,
			CvPath:         StoragePath,
			SourceFrom:     sourceFrom,
			CloudObject:    CloudObject,
		})

	} else if mode == "Update" {
		fmt.Println(mobile)
		err = candidates.C.Update(bson.M{"mobile": OldMobile}, &models.ApplicantInfo{
			Name:           name,
			Age:            age,
			Gender:         gender,
			Mobile:         mobile,
			Email:          email,
			Location:       location,
			Qualification:  qualification,
			Specialization: specialization,
			Department:     department,
			JobCode:        jobCode,
			Position:       position,
			Experience:     experience,
			CvPath:         StoragePath,
			SourceFrom:     sourceFrom,
			CloudObject:    CloudObject,
		})
	}

	if err != nil {
		panic(err)
	}
}
