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
		apiMessageResponse(w, http.StatusUnauthorized, "Authentication needed")
		return
	}
	ka.userID = userID
	switch {
	case r.Method == http.MethodPost && keyAllApi.MatchString(r.URL.Path):
		ka.createApiKeyHandler(w, r)
		return
	case r.Method == http.MethodDelete && keySingleApi.MatchString(r.URL.Path):
		ka.deleteApiKeyHandler(w, r)
		return
	case r.Method == http.MethodGet && keyAllApi.MatchString(r.URL.Path):
		ka.listApiKeysHandler(w, r)
		return
	default:
		apiMessageResponse(w, http.StatusNotFound, "not found")
		return
	}
}

func (ka *keyAPI) listApiKeysHandler(w http.ResponseWriter, _ *http.Request) {
	apiKeys, err := services.ListAPIKeys()
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(apiKeys)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (ka *keyAPI) createApiKeyHandler(w http.ResponseWriter, _ *http.Request) {
	apiKey, err := services.CreateAPIKey(ka.userID)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(apiKey)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (ka *keyAPI) deleteApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	matches := keySingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	apiKey := matches[1]
	_, err := services.DeleteAPIKey(apiKey, ka.userID)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	apiMessageResponse(w, http.StatusOK, "API key deleted")
}
