package web

import (
	"integrand/services"
	"log/slog"
	"net/http"
	"regexp"
)

var (
	keyAllApi    = regexp.MustCompile(`^\/api\/v1\/apikey[\/]*$`)
	keySingleApi = regexp.MustCompile(`^\/api\/v1\/apikey\/(.*)$`)
)

type keyAPI struct {
	userID int
}

func (ka *keyAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	userID, err := services.AuthenticateCookie(w, r)
	if err != nil {
		slog.Error(err.Error())
		// TODO: replace this with proper not authenticated response
		notFoundApiError(w)
		return
	}
	ka.userID = userID
	switch {
	case r.Method == http.MethodPost && keyAllApi.MatchString(r.URL.Path):
		ka.createApiKeyHandler(w, r)
	case r.Method == http.MethodDelete && keySingleApi.MatchString(r.URL.Path):
		ka.deleteApiKeyHandler(w, r)
	default:
		notFoundApiError(w)
	}
}

func (ka *keyAPI) createApiKeyHandler(w http.ResponseWriter, _ *http.Request) {
	apiKey, err := services.CreateAPIKey(ka.userID)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	response := map[string]string{"apiKey": apiKey}
	resJsonBytes, err := generateSuccessMessage(response)
	if err != nil {
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (ka *keyAPI) deleteApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	matches := keySingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFoundApiError(w)
		return
	}
	apiKey := matches[1]
	_, err := services.DeleteAPIKey(apiKey, ka.userID)
	if err != nil {
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	resJsonBytes, _ := generateSuccessMessage(map[string]string{"message": "API key deleted"})
	w.Write(resJsonBytes)
}
