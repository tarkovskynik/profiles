package main

import (
	"context"
	"net"

	"github.com/profiles/database"
	"github.com/profiles/server"
	"github.com/profiles/users"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs"
)

func main() {
	ctx := context.Background()
	config := server.Config{
		Address:       "127.0.0.1:8087",
		DbAddress:     "testdb:12345@tcp(localhost:3306)/dbtest?multiStatements=true",
		MigrationPath: "file:///Users/mykytatarkovskyi/apps/profiles/database/migration",
		DBName:        "dbtest",
	}

	db, err := database.NewDB(config)
	if err != nil {
		logrus.Fatalf("Error connecting to Database: %s", err.Error())
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

	userService := users.NewService(db.Users())

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		logrus.Fatalf("%s", err.Error())
	}

	server := server.NewServer(config, listener, userService)

	err = server.Run(ctx)
	if err != nil {
		logrus.Fatalf("%s", err.Error())
	}
}
