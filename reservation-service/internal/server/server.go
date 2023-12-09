package server

import (
	"bmstu-dips-lab2/config"
	"bmstu-dips-lab2/pkg/postgres"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	db     *postgres.Postgres
}

func New(cfg *config.Config, db *postgres.Postgres) *Server {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		// cors.New(cors.Config{
		// 	AllowOrigins: cfg.Cors.AllowOrigins,
		// 	AllowMethods: cfg.Cors.AllowMethods,
		// 	AllowHeaders: cfg.Cors.AllowHeaders,
		// }),
	)
	return &Server{router: router, cfg: cfg, db: db}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		Handler:        s.router,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	if err := s.MapHandlers(); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	return server.Shutdown(ctx)
}
