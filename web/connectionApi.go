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
		apiMessageResponse(w, http.StatusUnauthorized, "Authentication needed")
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
		apiMessageResponse(w, http.StatusNotFound, "not found")
		return
	}
}

func (ga *glueAPI) endpointHandler(w http.ResponseWriter, r *http.Request) {
	security_key := r.URL.Query().Get("apikey")
	// Check the content type header and parse appropriately
	matches := glueEndpointSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	// Check if we hit the rights one....
	sticky, err := services.GetEndpointBySecurityKey(matches[1], security_key)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	var i interface{}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	err = services.MessageToSink(sticky.TopicName, sticky.UserId, i)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	apiMessageResponse(w, http.StatusOK, "message sent successfully")
}

func (ga *glueAPI) getAllGlueHandlers(w http.ResponseWriter, _ *http.Request) {
	endpoints, err := services.GetEndpoints(ga.userID)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(endpoints)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (ga *glueAPI) getGlueHandler(w http.ResponseWriter, r *http.Request) {
	matches := glueSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	stickyConnection, err := services.GetEndpoint(matches[1], ga.userID)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(stickyConnection)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
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
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	stickyConnection, err := services.CreateEndpoint(ga.userID, createBody.RouteID, createBody.TopicName)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(stickyConnection)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (ga *glueAPI) deleteGlueHandler(w http.ResponseWriter, r *http.Request) {
	matches := glueSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	_, err := services.RemoveEndpoint(ga.userID, matches[1])
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	c := map[string]interface{}{"msg": "successfully deleted glue handler"}
	resJsonBytes, err := generateSuccessMessage(c)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}
