package Routers

import (
	"net/http"

	handlers "github.com/SivaShhankar/CMS_Cloud/Handlers"
	"github.com/codegangsta/negroni"

	"github.com/gorilla/mux"
)

// SetCandidateRoutes -- Candidate Handlers mapping
func SetCandidateRoutes(router *mux.Router) *mux.Router {
	taskRouter := mux.NewRouter()
	fs := http.FileServer(http.Dir("Templates"))
	router.PathPrefix("/css/").Handler(fs)
	router.PathPrefix("/images/").Handler(fs)
	router.PathPrefix("/JS/").Handler(fs)
	router.PathPrefix("/Files/").Handler(fs)

	taskRouter.Handle("/Index", http.HandlerFunc(handlers.Index))
	taskRouter.Handle("/Upload", http.HandlerFunc(handlers.Upload))
	taskRouter.Handle("/View", http.HandlerFunc(handlers.View))
	taskRouter.Handle("/Delete", http.HandlerFunc(handlers.Delete))
	taskRouter.Handle("/Search", http.HandlerFunc(handlers.Search))
	taskRouter.Handle("/Filter", http.HandlerFunc(handlers.Filter))
	taskRouter.Handle("/EditData", http.HandlerFunc(handlers.Edit))

	//[START AUTHORIZATION - MIDDLEWARE]
	router.PathPrefix("/Index").Handler(negroni.New(
		negroni.HandlerFunc(handlers.Authorize),
		negroni.Wrap(taskRouter),
	))
	router.PathPrefix("/Upload").Handler(negroni.New(
		negroni.HandlerFunc(handlers.Authorize),
		negroni.Wrap(taskRouter),
	))
	router.PathPrefix("/View").Handler(negroni.New(
		negroni.HandlerFunc(handlers.Authorize),
		negroni.Wrap(taskRouter),
	))
	router.PathPrefix("/Delete").Handler(negroni.New(
		negroni.HandlerFunc(handlers.Authorize),
		negroni.Wrap(taskRouter),
	))
	router.PathPrefix("/Search").Handler(negroni.New(
		negroni.HandlerFunc(handlers.Authorize),
		negroni.Wrap(taskRouter),
	))
	router.PathPrefix("/Filter").Handler(negroni.New(
		negroni.HandlerFunc(handlers.Authorize),
		negroni.Wrap(taskRouter),
	))
	router.PathPrefix("/EditData").Handler(negroni.New(
		negroni.HandlerFunc(handlers.Authorize),
		negroni.Wrap(taskRouter),
	))
	//[END AUTHORIZATION]
	return router
}
