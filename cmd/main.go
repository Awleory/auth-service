package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/awleory/medodstest/internal/config"
	"github.com/awleory/medodstest/internal/service"
	"github.com/awleory/medodstest/internal/storage/psql"
	"github.com/awleory/medodstest/internal/transport/rest"
	"github.com/awleory/medodstest/pkg/database"
	"github.com/awleory/medodstest/pkg/hash"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

const (
	CONFIG_DIR  = "config"
	CONFIG_FILE = "main"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connection(database.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	usersService := service.NewUsers(
		psql.NewUsers(db),
		psql.NewTokens(db),
		hash.NewSHA1Hasher(os.Getenv("SALT")),
		[]byte(os.Getenv("JWT_SECRET_KEY")))

	handler := rest.NewHandler(usersService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Info("SERVER STARTED")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
