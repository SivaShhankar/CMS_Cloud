package Handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"

	controllers "github.com/SivaShhankar/CMS_Cloud/Controllers"
	config "github.com/SivaShhankar/CMS_Cloud/Database"
)

type LoginInfo struct {
	LoginStatus string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}

	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := http.Cookie{
			Name:  "sessionUser",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, &cookie)
		fmt.Println("Cookie Added")
	}

}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("sessionUser"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func clearSession(response http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "sessionUser",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, &cookie)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	clearSession(w)
	http.Redirect(w, r, "/Login", http.StatusSeeOther)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		t, _ := template.ParseFiles("Templates/ResetPassword.html")
		t.Execute(w, nil)

	} else {

		r.ParseMultipartForm(32 << 20)

		newPassword := r.FormValue("NewPassword")

		if newPassword != "password" {

			userName := getUserName(r)
			error := controllers.ResetPassword(userName, newPassword)

			if error == nil {
				http.Redirect(w, r, "/Index", http.StatusSeeOther)
			} else {
				fmt.Println(error)
			}

		} else {

			t, _ := template.ParseFiles("Templates/ResetPassword.html")
			t.Execute(w, nil)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)

	if r.Method == "GET" {

		l := LoginInfo{}

		l.LoginStatus = ""

		t, _ := template.ParseFiles("Templates/Login.html")

		t.Execute(w, l)

	} else {

		r.ParseMultipartForm(32 << 20)

		userName := r.FormValue("Email")
		password := r.FormValue("Password")

		isCorrect, resetPasswordStatus := controllers.ValidateUser(config.Session, userName, password)

		if isCorrect {
			setSession(userName, w)

			if resetPasswordStatus {
				http.Redirect(w, r, "/ResetPassword", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/Index", http.StatusSeeOther)
			}

		} else {

			l := LoginInfo{}

			l.LoginStatus = "Invalid Username or Password"

			t, _ := template.ParseFiles("Templates/Login.html")

			t.Execute(w, l)
		}
	}
}
