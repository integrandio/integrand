package web

import (
	"context"
	"integrand/services"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func setUpGothSession() {
	//WTF why do we need all this junk?
	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	googleCallbackUrl := os.Getenv("GOOGLE_CALLBACK_URL")
	googleClientKey := os.Getenv("GOOGLE_CLIENT_KEY")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	githubCallbackUrl := os.Getenv("GITHUB_CALLBACK_URL")
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	goth.UseProviders(
		google.New(
			googleClientKey,
			googleClientSecret,
			googleCallbackUrl,
			"email",
			"profile"),
		github.New(
			githubClientID,
			githubClientSecret,
			githubCallbackUrl,
			"user",
		),
	)
}

func googleCallback(w http.ResponseWriter, r *http.Request) {
	//TODO: rewrite this so we don't have all the crap dependencies of gothic
	goUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	user, err := services.GoogleAuthenticate(goUser.Email, goUser.UserID)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}

	sess := services.GetSession(w, r)
	sess.Set("userID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func googleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := r.WithContext(context.WithValue(ctx, gothic.ProviderParamKey, "google"))
	gothic.BeginAuthHandler(w, req)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	//TODO: rewrite this so we don't have all the crap dependencies of gothic
	goUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	user, err := services.GithubAuthenticate(goUser.Email, goUser.UserID)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}

	sess := services.GetSession(w, r)
	sess.Set("userID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func githubAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := r.WithContext(context.WithValue(ctx, gothic.ProviderParamKey, "github"))
	gothic.BeginAuthHandler(w, req)
}
