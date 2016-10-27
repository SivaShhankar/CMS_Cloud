package Handlers

import (
	//"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	uuid "github.com/satori/go.uuid"

	"golang.org/x/net/context"

	controllers "github.com/SivaShhankar/CMS_Cloud/Controllers"
	config "github.com/SivaShhankar/CMS_Cloud/Database"
	models "github.com/SivaShhankar/CMS_Cloud/Models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	LoginUserInfo *models.UserInfo
	AppSession    *sessions.Session
	SessionStore  sessions.Store
	SessionID     string
)

type Profile struct {
	Uname string
}

var (
	// authKey = []byte(securecookie.GenerateRandomKey(32))
	// encKey  = []byte(securecookie.GenerateRandomKey(32))

	//store = sessions.NewCookieStore(authKey, encKey)

	GoogleOauthConfig = &oauth2.Config{

		RedirectURL:  "http://localhost:8080/GoogleCallback",
		ClientID:     "208027129669-01q79kp88k1roi53rguj9qluo0ce0np3.apps.googleusercontent.com",
		ClientSecret: "hPcD4-VgM3m-mjJWAC_hcGQl",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString string // = randToken()

)

const (
	htmlIndex = "<html><body><a href='/GoogleLogin'>Log in with Google</a></body></html>"
)

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

//[START SESSION PART]
func Init() {
	gob.Register(&oauth2.Token{})
	gob.Register(LoginUserInfo)
	SessionID = uuid.NewV4().String()
	//var err error
	var authKey = []byte(securecookie.GenerateRandomKey(32))
	//var encKey = []byte(securecookie.GenerateRandomKey(32))
	cookieStore := sessions.NewCookieStore([]byte(authKey))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	SessionStore = cookieStore
	fmt.Println("Seesion ID--", SessionID)
}

// func initSession(r *http.Request) *sessions.Session {
// 	gob.Register(&oauth2.Token{})
// 	log.Println("session before get", AppSession)

// 	if AppSession != nil {
// 		return AppSession
// 	}

// 	// session, err := store.Get(r, "mycmssession")
// 	// session.Options = &sessions.Options{
// 	// 	Path:     "/",
// 	// 	MaxAge:   -1,
// 	// 	HttpOnly: true,
// 	// 	Secure:   true,
// 	// 	//Domain:   "https://cmscloud-145306.appspot.com",
// 	// }

// 	AppSession = session

// 	log.Println("session after get", session)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return session
// }

//[END SESSION PART]

//[START OATUH PART]

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	// t, _ := template.ParseFiles("Templates/Login.html")
	// t.Execute(w, nil)
	http.Redirect(w, r, "/Login", http.StatusTemporaryRedirect)
}
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("Templates/Login.html")
	t.Execute(w, nil)
}
func HandleAccessDenied(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Sing Out")
	var ok bool
	if SessionID == "" {
		http.Redirect(w, r, "/SignOut", http.StatusFound)
		return
	}
	session, err := SessionStore.New(r, SessionID)
	if err != nil {
		http.Redirect(w, r, "/SignOut", http.StatusFound)
		return
	}
	LoginUserInfo, ok = session.Values["UName"].(*models.UserInfo)

	if !ok {
		http.Redirect(w, r, "/SignOut", http.StatusTemporaryRedirect)
	} else {
		session.Options.MaxAge = -1 // Clear session.
		err1 := session.Save(r, w)
		fmt.Println("error on Accessdenied", err1)
		t, _ := template.ParseFiles("Templates/AccessDenied.html")
		t.Execute(w, LoginUserInfo)
	}
	//AppSession = nil

}
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {

	Init()
	oauthStateString = randToken()
	url := GoogleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	fmt.Println(SessionID)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	// AppSession = nil
	// clearSession(w)
	fmt.Println("Sing Out")
	if SessionID == "" {
		http.Redirect(w, r, LogOutURL, http.StatusFound)
		return
	}
	session, err := SessionStore.New(r, SessionID)
	if err != nil {
		http.Redirect(w, r, LogOutURL, http.StatusFound)
		return
	}
	fmt.Println("Sing Out 1")
	session.Options.MaxAge = -1 // Clear session.
	err1 := session.Save(r, w)
	fmt.Println("Error on Clearing Session", err1)
	http.Redirect(w, r, LogOutURL, http.StatusTemporaryRedirect)

}
func HandleGoogleCallBack(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	if state != oauthStateString {
		url := GoogleOauthConfig.AuthCodeURL(oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	ctx := context.Background()
	code := r.FormValue("code")
	fmt.Println("CODE", code)
	token, err := GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Redirect(w, r, "/GoogleLogin", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println("Session saved ...")
	session, err := SessionStore.New(r, SessionID)
	//session := initSession(r)
	session.Values["mycmstoken"] = token

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	var data = new(models.UserInfo)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	json.Unmarshal(contents, data)
	// LoginUserInfo = new(models.UserInfo)
	// LoginUserInfo = data
	session.Values["UName"] = data
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isCorrect, _ := controllers.ValidateUser(config.Session, data.EMail, "")
	fmt.Println("validting User -", isCorrect, data.EMail, config.Session)

	if !isCorrect {
		// http.SetCookie(w, &http.Cookie{
		// 	Name:    "User",
		// 	Expires: expire,
		// 	Value:   LoginUserInfo.EMail,
		// })
		http.Redirect(w, r, "https://www.google.com/accounts/Logout?continue=https://appengine.google.com/_ah/logout?continue=http://localhost:8080/AccessDenied", http.StatusTemporaryRedirect)
		// AppSession = nil
		return
	}

	http.Redirect(w, r, "/Index", http.StatusFound)
}

//[END OAUTH PART]
