package users

import (
	"context"

	"github.com/zeebo/errs"
)

var (
	// ErrNoUser indicated that user does not exist.
	ErrNoUser = errs.Class("user does not exist")
)

// DB exposes access to users db.
type DB interface {
	// ListOfUserProfiles returns all user profiles from the database.
	ListOfUserProfiles(ctx context.Context) ([]UserProfile, error)
	// GetUserProfileByName returns user profile by name from the database.
	GetUserProfileByName(ctx context.Context, userName string) (UserProfile, error)
	// ListOfApiKeys returns api keys from the database.
	ListOfApiKeys(ctx context.Context) ([]Auth, error)
}

// User describes user entity.
type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"userName"`
}

// Profile for user profile.
type Profile struct {
	UserID    int64  `json:"-"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"-"`
	Address   string `json:"-"`
	City      string `json:"city"`
}

// UserData for user school information.
type UserData struct {
	UserID int64  `json:"-"`
	School string `json:"school"`
}

type UserProfile struct {
	User
	Profile
	UserData
}

// Auth for authorized authentication.
type Auth struct {
	ID     int64  `json:"id"`
	ApiKey string `json:"apiKey"`
}
