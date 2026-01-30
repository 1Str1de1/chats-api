package server

import (
	"chats-api/internal/config"
	"chats-api/internal/handler"
	"chats-api/internal/model"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

type Server struct {
	router  http.Handler
	conf    *config.Config
	logger  *slog.Logger
	handler *handler.Handler
}

func NewServer(conf *config.Config) (*Server, error) {
	logger, err := setupLogger()
	if err != nil {
		return nil, errors.New("logger error: " + err.Error())
	}

	db, err := model.NewDB(conf.PostgresConf)
	if err != nil {
		return nil, errors.New("db error: " + err.Error())
	}

	if err := runMigrations(db); err != nil {
		return nil, errors.New("error " + err.Error())
	}

	h := handler.NewHandler(db, logger)

	hdlr := configureMux(h, conf.ApiVersion)

	return &Server{
		router:  hdlr,
		logger:  logger,
		conf:    conf,
		handler: h,
	}, nil
}

func (s *Server) Start() error {
	s.logger.Info("starting server...")

	if err := http.ListenAndServe("0.0.0.0:8080", s.router); err != nil {
		return err
	}
	return nil
}

func setupLogger() (*slog.Logger, error) {
	//projDir, err := os.Getwd()
	//if err != nil {
	//	return nil, err
	//}
	//
	//logFile := filepath.Join(projDir, "logs", "logs.log")
	//file, err := os.OpenFile(
	//	logFile,
	//	os.O_CREATE|os.O_WRONLY|os.O_APPEND,
	//	0666)
	//if err != nil {
	//	return nil, err
	//}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return logger, nil
}

func runMigrations(db *gorm.DB) error {

	sql, err := db.DB()
	if err != nil {
		return err
	}

	migrationsDir := "migrations"
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	goose.SetLogger(log.New(os.Stdout, "[goose]", 0))

	if err := goose.DownTo(sql, migrationsDir, 0); err != nil {
		return err
	}
	if err := goose.Up(sql, migrationsDir); err != nil {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}

func configureMux(h *handler.Handler, apiVersion string) http.Handler {
	mux := http.NewServeMux()

	apiPrefix := fmt.Sprintf("/api/%s/chats", apiVersion)

	mux.HandleFunc("POST "+apiPrefix, h.HandleChatsCreate())
	mux.HandleFunc("POST "+apiPrefix+"/{id}/messages", h.HandleMessagesCreate())
	mux.HandleFunc("GET "+apiPrefix+"/{id}", h.HandleMessagesGet())
	mux.HandleFunc("DELETE "+apiPrefix+"/{id}", h.HandleChatsDelete())

	return mux
}
