package web

import (
	"encoding/json"
	"integrand/services"
	"log/slog"
	"net/http"
	"regexp"
)

var (
	glueEndpointSingleApi = regexp.MustCompile(`^\/api/v1/connector/f\/(.*)$`)
	//TODO: create a route for the routes....
	glueAllApi    = regexp.MustCompile(`^\/api/v1/connector[\/]*$`)
	glueSingleApi = regexp.MustCompile(`^\/api/v1/connector\/(.*)$`)
)

type glueAPI struct {
	userID int
}

func (ga *glueAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case glueEndpointSingleApi.MatchString(r.URL.Path):
		enableCors(&w)
		ga.endpointHandler(w, r)
	default:
		ga.apier(w, r)
	}
}

func (ga *glueAPI) apier(w http.ResponseWriter, r *http.Request) {
	userId, err := apiBrowserAPIAuthenticate(w, r)
	if err != nil {
		slog.Error(err.Error())
		notFoundApiError(w)
		return
	}
	ga.userID = userId
	w.Header().Set("content-type", "application/json")
	switch {
	//Glue API Routes
	case r.Method == http.MethodGet && glueAllApi.MatchString(r.URL.Path):
		ga.getAllGlueHandlers(w, r)
		return
	case r.Method == http.MethodGet && glueSingleApi.MatchString(r.URL.Path):
		ga.getGlueHandler(w, r)
		return
	case r.Method == http.MethodPost && glueAllApi.MatchString(r.URL.Path):
		ga.createGlueHandler(w, r)
		return
	case r.Method == http.MethodDelete && glueSingleApi.MatchString(r.URL.Path):
		ga.deleteGlueHandler(w, r)
		return
	default:
		notFoundApiError(w)
		return
	}
}

func (ga *glueAPI) endpointHandler(w http.ResponseWriter, r *http.Request) {
	security_key := r.URL.Query().Get("apikey")
	// Check the content type header and parse appropriately
	matches := glueEndpointSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFoundApiError(w)
		return
	}
	// Check if we hit the rights one....
	sticky, err := services.GetStickyConnectionBySecurityKey(matches[1], security_key)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	var i interface{}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	err = services.MessageToSink(sticky.TopicName, sticky.UserId, i)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	c := map[string]interface{}{"msg": "message sent successfully"}
	resJsonBytes, err := generateSuccessMessage(c)
	if err != nil {
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (ga *glueAPI) getAllGlueHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		stickey_connections, err := services.GetStickyConnections(ga.userID)
		if err != nil {
			slog.Error(err.Error())
			internalServerError(w)
			return
		}
		resJsonBytes, err := generateSuccessMessage(stickey_connections)
		if err != nil {
			slog.Error(err.Error())
			internalServerError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resJsonBytes)
	default:
		notFoundApiError(w)
	}
}

func (ga *glueAPI) getGlueHandler(w http.ResponseWriter, r *http.Request) {
	matches := glueSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFoundApiError(w)
		return
	}
	stickyConnection, err := services.GetStickyConnection(ga.userID, matches[1])
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	resJsonBytes, err := generateSuccessMessage(stickyConnection)
	if err != nil {
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

type CreateGlueBody struct {
	RouteID   string `json:"id"`
	TopicName string `json:"topicName"`
}

func (ga *glueAPI) createGlueHandler(w http.ResponseWriter, r *http.Request) {
	var createBody CreateGlueBody
	if err := json.NewDecoder(r.Body).Decode(&createBody); err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	stickyConnection, err := services.CreateStickyConnection(ga.userID, createBody.RouteID, createBody.TopicName)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	resJsonBytes, err := generateSuccessMessage(stickyConnection)
	if err != nil {
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (ga *glueAPI) deleteGlueHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		matches := glueSingleApi.FindStringSubmatch(r.URL.Path)
		if len(matches) < 2 {
			notFoundApiError(w)
			return
		}
		_, err := services.RemoveStickyConnection(ga.userID, matches[1])
		if err != nil {
			slog.Error(err.Error())
			internalServerError(w)
			return
		}
		c := map[string]interface{}{"msg": "successfully deleted glue handler"}
		resJsonBytes, err := generateSuccessMessage(c)
		if err != nil {
			internalServerError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resJsonBytes)
	default:
		notFoundApiError(w)
	}
}
