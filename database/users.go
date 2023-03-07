package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/profiles/users"
	"github.com/zeebo/errs"
)

// ErrUsers indicates that there was an error in the database.
var ErrUsers = errs.Class("users repository error")

// usersDB provides access to users db.
type usersDB struct {
	conn *sql.DB
}

// ListOfUserProfiles returns all user profiles from the database.
func (usersDB *usersDB) ListOfUserProfiles(ctx context.Context) ([]users.UserProfile, error) {
	query := `SELECT id, username, first_name, last_name, city, school 
				FROM user 
				LEFT JOIN user_profile ON user.id = user_profile.user_id
				LEFT JOIN user_data ON user.id = user_data.user_id`

	rows, err := usersDB.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrUsers.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	var data []users.UserProfile
	for rows.Next() {
		var userProfile users.UserProfile
		err := rows.Scan(&userProfile.ID, &userProfile.UserName, &userProfile.FirstName, &userProfile.LastName, &userProfile.City, &userProfile.School)
		if err != nil {
			return nil, ErrUsers.Wrap(err)
		}

		data = append(data, userProfile)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrUsers.Wrap(err)
	}

	return data, ErrUsers.Wrap(err)
}

// GetUserProfileByName returns user profile by name from the database.
func (usersDB *usersDB) GetUserProfileByName(ctx context.Context, userName string) (users.UserProfile, error) {
	var userProfile users.UserProfile
	query := `SELECT id, username, first_name, last_name, city, school 
				FROM user 
				LEFT JOIN user_profile ON user.id = user_profile.user_id
				LEFT JOIN user_data ON user.id = user_data.user_id
				WHERE username = ?`

	row := usersDB.conn.QueryRowContext(ctx, query, userName)

	err := row.Scan(&userProfile.ID, &userProfile.UserName, &userProfile.FirstName, &userProfile.LastName, &userProfile.City, &userProfile.School)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users.UserProfile{}, users.ErrNoUser.Wrap(err)
		}
		return users.UserProfile{}, ErrUsers.Wrap(err)
	}

	return userProfile, nil
}

// ListOfApiKeys returns api keys from the database.
func (usersDB *usersDB) ListOfApiKeys(ctx context.Context) ([]users.Auth, error) {
	query := `SELECT * FROM auth`

	rows, err := usersDB.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrUsers.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	var data []users.Auth
	for rows.Next() {
		var auth users.Auth
		err := rows.Scan(&auth.ID, &auth.ApiKey)
		if err != nil {
			return nil, ErrUsers.Wrap(err)
		}

		data = append(data, auth)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrUsers.Wrap(err)
	}

	return data, ErrUsers.Wrap(err)
}
