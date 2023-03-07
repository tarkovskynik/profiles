package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/profiles/users"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs"
)

var (
	// ErrUsers is an internal error type for users controller.
	ErrUsers = errs.Class("users controller error")
)

// Users is a mvc controller that handles all users related views.
type Users struct {
	users *users.Service
}

// NewUsers is a constructor for users controller.
func NewUsers(users *users.Service) *Users {
	usersController := &Users{
		users: users,
	}

	return usersController
}

// GetProfile is an endpoint that returns the current user profile with all relevant information.
func (controller *Users) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var result interface{}
	var err error
	query := r.URL.Query()
	userName := query.Get("username")

	if len(userName) > 0 {
		result, err = controller.users.GetUserProfileByUserName(ctx, userName)
		if err != nil {
			logrus.Errorf("could not get profile: %v", ErrUsers.Wrap(err))
			controller.serveError(w, http.StatusInternalServerError, ErrUsers.Wrap(err))
			return
		}
	} else {
		result, err = controller.users.UserProfiles(ctx)
		if err != nil {
			logrus.Errorf("could not get profile: %v", ErrUsers.Wrap(err))
			controller.serveError(w, http.StatusInternalServerError, ErrUsers.Wrap(err))
			return
		}
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		logrus.Errorf("failed to write json response: %v", ErrUsers.Wrap(err))
		return
	}
}

// serveError replies to request with specific code and error.
func (controller *Users) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		logrus.Errorf("failed to write json error response: %v", ErrUsers.Wrap(err))
	}
}
