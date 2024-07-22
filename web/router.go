package web

import (
	"encoding/json"
	"integrand/services"
	"log"
	"net/http"
)

func NewNewWebRouter() *http.ServeMux {
	mux := http.NewServeMux()
	//Authentication stuff
	mux.HandleFunc("/login", LoginPage)
	// mux.HandleFunc("/register", RegisterPage)

	// Application UI
	mux.HandleFunc("/", applicationPage)

	glueApi := &glueAPI{}
	mux.Handle("/api/v1/connector", glueApi)
	mux.Handle("/api/v1/connector/", glueApi)

	topicApi := &topicAPI{}
	mux.Handle("/api/v1/topic", topicApi)
	mux.Handle("/api/v1/topic/", topicApi)

	// Register keyAPI
	keyApi := &keyAPI{}
	mux.Handle("/api/v1/apikey", keyApi)
	mux.Handle("/api/v1/apikey/", keyApi)

	// Serve static files from the "static" directory.
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	return mux
}

func apiBrowserAPIAuthenticate(w http.ResponseWriter, r *http.Request) (int, error) {
	// First let's try to authenticate with our session
	userID, err := services.AuthenticateCookie(w, r)
	if err != nil {
		// Enable coors and then check for token
		enableCors(&w)
		// This is publicly exposed, we need to protect with a token
		authorizationHeader := r.Header.Get("Authorization")
		apiKey, err := services.AuthenticateToken(authorizationHeader)
		return apiKey.UserId, err
	}

	return userID, nil
}

func sessionAuthenticateOrRedirect(w http.ResponseWriter, r *http.Request) {
	_, err := services.AuthenticateCookie(w, r)
	if err != nil {
		// Successful login, redirect to a welcome page.
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Expose-Headers", "Content-Type")
}

func apiMessageResponse(w http.ResponseWriter, errorCode int, message string) {
	body := map[string]string{
		"message": message,
	}
	// This should never error...
	bodyBytes, _ := json.Marshal(body)

	w.WriteHeader(errorCode)
	w.Write(bodyBytes)
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
