package web

import (
	"encoding/json"
	"integrand/services"
	"log"
	"net/http"
)

func NewNewWebRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", LoginPage)
	// Application UI
	mux.HandleFunc("/app", applicationPage)
	mux.HandleFunc("/app/", applicationPage)

	glueApi := &glueAPI{}
	mux.Handle("/api/v1/glue", glueApi)
	mux.Handle("/api/v1/glue/", glueApi)
	mux.Handle("/api/v1/glue/endpoint", glueApi)

	topicApi := &topicAPI{}
	mux.Handle("/api/v1/topic", topicApi)
	mux.Handle("/api/v1/topic/", topicApi)

	// Serve static files from the "static" directory.
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	return mux
}

func apiBrowserAPIAuthenticate(w http.ResponseWriter, r *http.Request) error {
	// First let's try to authenticate with our session
	sess := services.GetSession(w, r)
	user, err := sess.Get("email")
	if err != nil {
		log.Fatalln(err)
	}
	if user == nil {
		// Enable coors and then check for token
		enableCors(&w)
		// This is publicly exposed, we need to protect with a token
		header := r.Header.Get("Authorization")
		err := services.AuthorizeToken(header)
		return err
	}
	return nil
}

func sessionAuthenticate(w http.ResponseWriter, r *http.Request) {
	sess := services.GetSession(w, r)
	user, err := sess.Get("email")
	if err != nil {
		log.Fatalln(err)
	}
	if user == nil {
		// Successful login, redirect to a welcome page.
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func internalServerError(w http.ResponseWriter) {
	resBytes, _ := generateErrorMessage("internal server error")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resBytes)
}

func notFoundApiError(w http.ResponseWriter) {
	c := map[string]interface{}{"api": "not found"}
	resBytes, _ := generateFailMessage(c)
	w.WriteHeader(http.StatusNotFound)
	w.Write(resBytes)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Expose-Headers", "Content-Type")
}

// Apis will be formatted using the "jsend" spec

type apiResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func generateSuccessMessage(data interface{}) ([]byte, error) {
	res := apiResponse{
		Status: "success",
		Data:   data,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return resBytes, err
	}
	return resBytes, nil
}

func generateFailMessage(data interface{}) ([]byte, error) {
	res := apiResponse{
		Status: "fail",
		Data:   data,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return resBytes, err
	}
	return resBytes, nil
}

func generateErrorMessage(message string) ([]byte, error) {
	res := apiResponse{
		Status:  "error",
		Message: message,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return resBytes, err
	}
	return resBytes, nil
}
