package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"zion/internal/config"
	zerr "zion/internal/errors"
	"zion/internal/hash"
	"zion/internal/storage"
	"zion/internal/storage/db"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ZionServer struct {
	config        *config.Config
	db            *gorm.DB
	users         storage.UserStorageInterface
	sessions      storage.SessionStorageInterface
	todos         storage.TodoStorageInterface
	httpServer    *http.Server
	hash          *hash.PasswordHash
	sessionCookie string
	router        chi.Router
	wg            sync.WaitGroup
}

func NewZionServer(
	cfg *config.Config,
	dbConn *gorm.DB,
	userStorage storage.UserStorageInterface,
	sessionStorage storage.SessionStorageInterface,
	todoStorage storage.TodoStorageInterface,
	hash *hash.PasswordHash,
) *ZionServer {
	s := &ZionServer{
		config:        cfg,
		db:            dbConn,
		users:         userStorage,
		sessions:      sessionStorage,
		todos:         todoStorage,
		hash:          hash,
		sessionCookie: cfg.SessionCookieName,
		httpServer: &http.Server{
			Addr:        fmt.Sprintf(":%s", cfg.Port),
			IdleTimeout: time.Minute,
			// ReadTimeout:    10 * time.Second,
			// WriteTimeout:   30 * time.Second,
			// MaxHeaderBytes: 1 << 20,
		},
	}
	s.SetupRouter()
	s.httpServer.Handler = s.router
	return s
}

func InitializeZionServer(cfg *config.Config) (*ZionServer, error) {
	// (1) Create password hash
	passwordHash := hash.NewPasswordHash()

	// (2) Connect to the database
	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		return nil, zerr.ErrFailedToConnectToDB
	}

	// (3) Create the database models (or migrate)
	err = db.CreateModels(dbConn, cfg.DatabaseMode, passwordHash)
	if err != nil {
		return nil, zerr.ErrCreateTables
	}

	// (4) Create storage layers
	userStorage := storage.NewUserStorage(storage.UserStorageParams{
		DB:           dbConn,
		PasswordHash: passwordHash,
	})
	sessionStorage := storage.NewSessionStorage(storage.SessionStorageParams{
		DB: dbConn,
	})
	todoStorage := storage.NewTodoStorage(storage.TodoStorageParams{
		DB: dbConn,
	})

	// (5) Create Zion Server
	server := NewZionServer(
		cfg,
		dbConn,
		userStorage,
		sessionStorage,
		todoStorage,
		passwordHash,
	)

	return server, nil
}

func (s *ZionServer) Start() {
	// Initialize kill signals
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	// Increment the WaitGroup counter for the HTTP server goroutine
	s.wg.Add(1)
	go func() {
		defer s.wg.Done() // Decrement the WaitGroup counter when the server stops
		err := s.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("server is shutting down...")
		} else if err != nil {
			log.Fatalf("server error: %v.", err)
		}
	}()
	log.Printf("starting server on port %s", s.config.Port)

	// Block until a kill signal is received
	<-killSig

	// Create a shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	} else {
		log.Print("server gracefully shutdown.")
	}

	// Wait for all background goroutines to finish
	s.wg.Wait()
	log.Print("all background operations complete. server shutdown complete.")
}

func (s *ZionServer) SetupRouter() {
	// Using chi router (maybe just use standard library in the future?)
	s.router = chi.NewRouter()

	// Initialize file server
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	s.router.Handle("/static/*", fs)
	s.router.Handle("/favicon.ico", fs)

	// Define routes
	s.EstablishRoutes()
}
