package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	go func() {
		if err := server.Run(ctx); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("app shutting down")

	if err := server.Close(); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
