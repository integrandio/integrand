package web

import (
	"fmt"
	"html/template"
	"integrand/services"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	sess := services.GetSession(w, r)
	switch r.Method {
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := services.EmailAuthenticate(email, password)
		if err != nil {
			log.Println(err)
			// Invalid credentials, show the login page with an error message.
			fmt.Fprintf(w, "Invalid credentials. Please try again.")
			return
		} else {
			err = sess.Set("userID", user.ID)
			if err != nil {
				log.Fatal(err)
			}
			// Successful login, redirect to a welcome page.
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	case http.MethodGet:
		user, err := sess.Get("userID")
		if err != nil {
			log.Fatal(err)
			return
		}
		if user != nil {
			// Successful login, redirect to a welcome page.
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// If not a POST request, serve the login page template.
		tmpl, err := template.ParseFiles("web/templates/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	default:
		notFoundApiError(w)
	}
}

func applicationPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sessionAuthenticateOrRedirect(w, r)
		fileContents, err := os.ReadFile("web/templates/baseApp.html")
		if err != nil {
			slog.Error(err.Error())
			internalServerError(w)
			return
		}
		templateString := strings.Replace(string(fileContents), "#replace_me#", "{{ template \"application.html\" }}", 1)
		tmpl, err := template.New("myTemplate").Parse(templateString)
		if err != nil {
			slog.Error(err.Error())
			internalServerError(w)
			return
		}
		_, err = tmpl.ParseFiles(
			"web/templates/appShell.html",
			"web/templates/application.html",
		)
		if err != nil {
			slog.Error(err.Error())
			internalServerError(w)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err)
			internalServerError(w)
			return
		}
	default:
		notFoundApiError(w)
	}
}
