package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zion/internal/config"
	"zion/internal/hash"
	"zion/internal/storage"
	"zion/internal/storage/db"

	"gorm.io/gorm"
)

type Server struct {
	port          string
	db            *gorm.DB
	users         *storage.UserStorage
	sessions      *storage.SessionStorage
	http          *http.Server
	hash          *hash.PasswordHash
	sessionCookie string
}

func CreateServer(cfg *config.Config) *Server {
	// Connect to the database
	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to the database: %v", err)
	}

	// Create the database models
	db.CreateModels(dbConn)

	// Initialize the server
	server := &Server{
		port: cfg.Port,
		http: &http.Server{
			Addr:           fmt.Sprintf(":%s", cfg.Port),
			IdleTimeout:    time.Minute,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		db:            dbConn,
		users:         storage.NewUserStorage(storage.UserStorageParameters{DB: dbConn, PasswordHash: ""}),
		sessions:      storage.NewSessionStorage(storage.SessionStorageParameters{DB: dbConn}),
		hash:          hash.NewPasswordHash(),
		sessionCookie: cfg.SessionCookieName,
	}

	// Set the server handler
	server.http.Handler = CreateRouter(server)

	return server
}

func (s *Server) Start() {
	// Initialize kill signals
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server
	go func() {
		err := s.http.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("❎ Server is shutting down...")
		} else if err != nil {
			log.Fatalf("❌ Server error: %v.", err)
		}
	}()
	log.Printf("🚀 Starting server on %s", s.http.Addr)
	<-killSig

	// Create a shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := s.http.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}
	log.Print("❎ Server shutdown complete.")
}
