package Routers

import (
	"net/http"

	"github.com/gorilla/mux"

	handlers "github.com/SivaShhankar/CMS_Cloud/Handlers"
)

// SetOAuthRoutes - OAuth Handler mapping
func SetOAuthRoutes(router *mux.Router) *mux.Router {

	//[OAuth Handlers Start]
	router.Handle("/", http.HandlerFunc(handlers.HandleIndex))                  //	Google Sign in
	router.Handle("/Login", http.HandlerFunc(handlers.HandleLogin))             //	Google Sign in
	router.Handle("/GoogleLogin", http.HandlerFunc(handlers.HandleGoogleLogin)) //
	router.Handle("/GoogleCallback", http.HandlerFunc(handlers.HandleGoogleCallBack))
	router.Handle("/SignOut", http.HandlerFunc(handlers.SignOut))
	router.Handle("/AccessDenied", http.HandlerFunc(handlers.HandleAccessDenied))

	//[OAuth Handlers End]

	return router
}
