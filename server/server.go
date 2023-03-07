package server

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profiles/server/controllers"
	"github.com/profiles/users"
	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"
)

var (
	// Error is an error class that indicates internal http server error.
	Error = errs.Class("console web server error")
)

// Config contains configuration for console web server.
type Config struct {
	Address       string
	DbAddress     string
	MigrationPath string
	DBName        string
}

// Server represents console web server.
type Server struct {
	config Config

	listener net.Listener
	server   http.Server

	users *users.Service
}

// NewServer is a constructor for console web server.
func NewServer(config Config, listener net.Listener, users *users.Service) *Server {
	server := &Server{
		config:   config,
		listener: listener,
		users:    users,
	}

	userController := controllers.NewUsers(users)

	router := mux.NewRouter()

	profileRouter := router.PathPrefix("/profile").Subrouter()
	profileRouter.Use(server.withAuth)
	profileRouter.HandleFunc("", userController.GetProfile).Methods(http.MethodGet)

	router.PathPrefix("/").HandlerFunc(server.appHandler)

	server.server = http.Server{
		Handler: router,
	}

	return server
}

// Run starts the server that host webapp and api endpoint.
func (server *Server) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()
		return Error.Wrap(server.server.Shutdown(context.Background()))
	})
	group.Go(func() error {
		defer cancel()
		err := server.server.Serve(server.listener)
		isCancelled := errs.IsFunc(err, func(err error) bool { return errors.Is(err, context.Canceled) })
		if isCancelled || errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		return Error.Wrap(err)
	})

	return Error.Wrap(group.Wait())
}

// Close closes server and underlying listener.
func (server *Server) Close() error {
	return Error.Wrap(server.server.Close())
}

// appHandler is web app http handler function.
func (server *Server) appHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=UTF-8")
	header.Set("Referrer-Policy", "same-origin")
}

// withAuth performs initial authorization before every request.
func (server *Server) withAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		apiKey := r.Header.Get("Api-key")

		err := server.users.Authentication(ctx, apiKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		handler.ServeHTTP(w, r.Clone(ctx))
	})
}
