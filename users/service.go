package users

import (
	"context"

	"github.com/zeebo/errs"
)

// ErrUsers indicates that there was an error in the service.
var ErrUsers = errs.Class("users service error")

// ErrUnauthenticated should be returned when user performs unauthenticated action.
var ErrUnauthenticated = errs.Class("user unauthenticated error")

// Service is handling users related logic.
type Service struct {
	users DB
}

// NewService is a constructor for users service.
func NewService(users DB) *Service {
	return &Service{
		users: users,
	}
}

// UserProfiles returns users profiles.
func (service *Service) UserProfiles(ctx context.Context) ([]UserProfile, error) {
	userProfiles, err := service.users.ListOfUserProfiles(ctx)
	if err != nil {
		return nil, ErrUsers.Wrap(err)
	}
	return userProfiles, nil
}

// GetUserProfileByUserName returns user by username.
func (service *Service) GetUserProfileByUserName(ctx context.Context, userName string) (UserProfile, error) {
	userProfile, err := service.users.GetUserProfileByName(ctx, userName)
	if err != nil {
		return UserProfile{}, ErrUsers.Wrap(err)
	}
	return userProfile, nil
}

// Authentication middleware check api key.
func (service *Service) Authentication(ctx context.Context, apiKey string) error {
	apiKeys, err := service.users.ListOfApiKeys(ctx)
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	for _, key := range apiKeys {
		if apiKey == key.ApiKey {
			return nil
		}
	}

	return ErrUsers.New("authentication wrong")
}
