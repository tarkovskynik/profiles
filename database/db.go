package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zeebo/errs"

	"github.com/profiles/server"
	"github.com/profiles/users"
)

var (
	Error = errs.Class("db error")
)

type DB interface {
	Users() users.DB
	Close() error
}

type database struct {
	conn *sql.DB
}

func NewDB(config server.Config) (DB, error) {
	conn, err := sql.Open("mysql", config.DbAddress)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	driver, _ := mysql.WithInstance(conn, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		config.MigrationPath,
		config.DBName,
		driver,
	)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		return nil, err
	}

	return &database{conn: conn}, nil
}

// Close closes underlying db connection.
func (db *database) Close() error {
	return Error.Wrap(db.conn.Close())
}

// Users provides access to accounts db.
func (db *database) Users() users.DB {
	return &usersDB{conn: db.conn}
}
