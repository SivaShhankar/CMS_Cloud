package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	cloudStorage "github.com/SivaShhankar/CMS_Cloud/CloudStorage"
	controllers "github.com/SivaShhankar/CMS_Cloud/Controllers"
	config "github.com/SivaShhankar/CMS_Cloud/Database"
	routers "github.com/SivaShhankar/CMS_Cloud/Routers"
)

func main() {

	// Load the configuration stuffs.
	config.LoadAppConfig()

	// Initiate the database information.
	config.CreateDBSession()

	// Add the neccessary indexes.
	config.AddIndexes()

	config.AddCredentialIndexes()

	controllers.CreateDefaultUserCredentials()

	// Configure Buckets
	cloudStorage.Init()

	// Configure Sessions
	//handlers.Init()

	// Created the routes of this application
	mux := mux.NewRouter().StrictSlash(false) //http.NewServeMux()
	mux = routers.SetCandidateRoutes(mux)
	//mux = routers.SetUserRoutes(mux)
	mux = routers.SetOAuthRoutes(mux)
	n := negroni.Classic()
	n.UseHandler(mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}

	log.Println("Listening in port:8080")

	// Listen the server.
	server.ListenAndServe()
	//http.ListenAndServe(":8080", mux)
}
