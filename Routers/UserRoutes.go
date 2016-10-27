package Routers

import (
	"net/http"

	"github.com/gorilla/mux"

	handlers "github.com/SivaShhankar/CMS_Cloud/Handlers"
)

func SetUserRoutes(router *mux.Router) *mux.Router {

	//router.Handle("/Login", http.HandlerFunc(handlers.Login))
	//router.Handle("/Logout", http.HandlerFunc(handlers.Logout))
	router.Handle("/ResetPassword", http.HandlerFunc(handlers.ResetPassword))

	return router
}
