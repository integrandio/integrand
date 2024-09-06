package web

import (
	"encoding/json"
	"integrand/services"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
)

var (
	userAllApi    = regexp.MustCompile(`^\/api\/v1\/user[\/]*$`)
	userSingleApi = regexp.MustCompile(`^\/api\/v1\/user\/(.*)$`)
)

type userAPI struct {
	userID int
}

type createUserBody struct {
	Email    string
	Password string
}

type updateUserBody struct {
	OldPassword string
	NewPassword string
}

func (u *userAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := apiBrowserAPIAuthenticate(w, r)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusUnauthorized, "Authentication needed")
		return
	}
	u.userID = userID
	switch {
	case r.Method == http.MethodGet && userAllApi.MatchString(r.URL.Path):
		u.getUsers(w, r)
		return
	case r.Method == http.MethodGet && userSingleApi.MatchString(r.URL.Path):
		u.getUser(w, r)
		return
	case r.Method == http.MethodPost && userAllApi.MatchString(r.URL.Path):
		u.createUser(w, r)
		return
	case r.Method == http.MethodPut && userSingleApi.MatchString(r.URL.Path):
		u.updateUser(w, r)
		return
	case r.Method == http.MethodDelete && userSingleApi.MatchString(r.URL.Path):
		u.deleteUser(w, r)
		return
	default:
		apiMessageResponse(w, http.StatusNotFound, "not found")
		return
	}
}

func (u *userAPI) getUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := services.GetUsers(u.userID)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(users)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (u *userAPI) getUser(w http.ResponseWriter, r *http.Request) {
	matches := userSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	user, err := services.GetUser(u.userID, id)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(user)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (u *userAPI) createUser(w http.ResponseWriter, r *http.Request) {
	var newUser createUserBody
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		apiMessageResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := services.CreateUser(u.userID, newUser.Email, newUser.Password)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(user)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (u *userAPI) updateUser(w http.ResponseWriter, r *http.Request) {
	matches := userSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	var updatedUser updateUserBody
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		apiMessageResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := services.UpdatePassword(u.userID, id, updatedUser.OldPassword, updatedUser.NewPassword)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(user)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (u *userAPI) deleteUser(w http.ResponseWriter, r *http.Request) {
	matches := userSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	user, err := services.RemoveUser(u.userID, id)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(user)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}
