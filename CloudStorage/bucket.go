package GCloudStorage

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	gstorage "google.golang.org/api/storage/v1"
)

const (
	// This scope allows the application full control over resources in Google Cloud Storage
	scope = gstorage.DevstorageFullControlScope
)

var (
	projectID  = "cmscloud-145306"
	bucketName = "cmscloud-145306.appspot.com"
)

//StorageService - variable for gcloud storage
var StorageService *gstorage.Service

// Init - Initialize google cloud storage
func Init() {
	var err1 error
	// Initialize context
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
	}
	StorageService, err1 = gstorage.New(client)
	fmt.Println(StorageService)
	if err1 != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	// Configure Buckets
	if _, err := StorageService.Buckets.Get(bucketName).Do(); err == nil {
		//fmt.Printf("Bucket %s already exists - skipping buckets.insert call.", bucketName)
	} else {
		// Create a bucket.
		if res, err := StorageService.Buckets.Insert(projectID, &gstorage.Bucket{Name: bucketName}).Do(); err == nil {
			fmt.Printf("Created bucket %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			fmt.Println(err)
		}
	}
}

// GCloudUploadFiles - Upload files into bucket
func GCloudUploadFiles(r *http.Request, fileName string) (string, string) {
	if StorageService == nil {
			return "",""
	}

	var extension = filepath.Ext(fileName)
	var metaData string

	if extension == ".docx" {
		metaData = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	} else if extension == ".doc" {
		metaData = "application/msword"
	} else if extension == ".pdf" {
		metaData = "application/pdf"
	}

	object := &gstorage.Object{Name: strings.Replace(fileName, " ", "-", -1), ContentType: metaData} // Map FileName as object name
	f, _, err := r.FormFile("file")                                                                  // Selected file from browser

	var fileMediaLink string
	fmt.Println(StorageService)

	//w: = storageService.Objects.Insert(bucketName,object).Media()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error Exited")
		return "", fileMediaLink
	}

	if res, err := StorageService.Objects.Insert(bucketName, object).Media(f).Do(); err == nil {
		fmt.Printf("File Created %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		fmt.Println("Failed to upload", err)
	}

	if res, err := StorageService.Objects.Get(bucketName, strings.Replace(fileName, " ", "-", -1)).Do(); err == nil {
		fileMediaLink = res.SelfLink // Storage path of uploaded file
	}
	return fileMediaLink, strings.Replace(fileName, " ", "-", -1)

}

// GCloudDeleteFiles - Delete existing object from bucket while in update mode
func GCloudDeleteFiles(file string) {

	if err := StorageService.Objects.Delete(bucketName, file).Do(); err == nil {
		fmt.Printf("Successfully deleted %s/%s during cleanup.\n\n", bucketName, file)
	} else {
		// If the object exists but wasn't deleted, the bucket deletion will also fail.
		fmt.Printf("Could not delete object during cleanup: %v\n\n", err)
	}
}
