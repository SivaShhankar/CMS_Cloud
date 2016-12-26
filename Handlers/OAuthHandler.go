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
	"time"

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

		RedirectURL:  "https://cmscloud-145306.appspot.com/GoogleCallback",
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

	// sessionID := GetCookieValue("ID", r)
	// Login_User = GetCookieValue("User", r)
	// if Login_User != "" {
	// 	fmt.Println("Session ID For Current Process", sessionID)
	// 	Login_User = GetCookieValue("User", r)

	// 	if Login_User == "" {
	// 		http.Redirect(w, r, LogOutURL, http.StatusFound)
	// 		return
	// 	}

	type Info struct {
		CurrentUser string
	}
	d := Info{CurrentUser: "Dummy"} //Login_User
	session, err := SessionStore.Get(r, SessionID)
	if err == nil {
		session.Options.MaxAge = -1 // Clear session.
		err1 := session.Save(r, w)
		fmt.Println("error on Accessdenied", err1)
	}
	clearCurrentSession("User", w)
	clearCurrentSession("ID", w)
	t, _ := template.ParseFiles("Templates/AccessDenied.html")
	t.Execute(w, d)
	//fmt.Println("Logged User", Login_User)
	// } else {
	// 	http.Redirect(w, r, LogOutURL, http.StatusFound)
	// 	return
	// }

}
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {

	Init()
	oauthStateString = randToken()
	url := GoogleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	fmt.Println(SessionID)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	sessionID := GetCookieValue("ID", r)
	if sessionID == "" {
		http.Redirect(w, r, LogOutURL, http.StatusFound)
		return
	}
	session, err := SessionStore.Get(r, sessionID)
	if err != nil {
		http.Redirect(w, r, LogOutURL, http.StatusFound)
		return
	}
	fmt.Println("Sing Out 1")
	session.Options.MaxAge = -1 // Clear session.
	err1 := session.Save(r, w)
	fmt.Println("Error on Clearing Session", err1)
	clearCurrentSession("User", w)
	clearCurrentSession("ID", w)
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

	setCookieValue("User", data.Name, w)
	setCookieValue("ID", SessionID, w)
	if !isCorrect {
		http.Redirect(w, r, "https://www.google.com/accounts/Logout?continue=https://appengine.google.com/_ah/logout?continue=https://cmscloud-145306.appspot.com/AccessDenied", http.StatusTemporaryRedirect)
		// AppSession = nil
		return
	}

	http.Redirect(w, r, "/Index", http.StatusFound)
}

func setCookieValue(ParamName string, ParamValue string, response http.ResponseWriter) {
	value := map[string]string{
		"name": ParamValue,
	}

	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := http.Cookie{
			Name:    ParamName,
			Value:   encoded,
			Path:    "/",
			Expires: time.Now().Add(356 * 24 * time.Hour),
		}
		http.SetCookie(response, &cookie)
		fmt.Println("Cookie Added")
	}

}

func GetCookieValue(ParamName string, request *http.Request) (ParamValue string) {
	if cookie, err := request.Cookie(ParamName); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			ParamValue = cookieValue["name"]
		}
	}
	return ParamValue
}

func clearCurrentSession(ParamName string, response http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   ParamName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, &cookie)
}

//[END OAUTH PART]
