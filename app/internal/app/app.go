package app

import (
	"LamodaTestTask/app/api/routes"
	"LamodaTestTask/app/pkg/logging"
	"LamodaTestTask/app/pkg/postgres"
	"context"
	"database/sql"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type App struct {
	//cfg        *config.Config
	router     *gin.Engine
	httpServer *http.Server
	pgClient   *sql.DB
}

func NewApp(ctx context.Context) (*App, error) {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := postgres.NewPostgresDB(&postgres.Config{
		Username: "postgres",
		Password: os.Getenv("POSTGRES_PASS"),
		Host:     "localhost",
		Port:     "5439",
		Database: "lamoda_db",
		SSLMode:  "disable",
	})

	router := routes.NewRouter(db)
	logging.GetLogger(ctx).Info("router initializing")

	return &App{
		router:   router,
		pgClient: db,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	logging.GetLogger(ctx).Info("application initialized and started")
	defer func() {
		if err := a.pgClient.Close(); err != nil {
			logging.GetLogger(ctx).Error(err)
		}
	}()

	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return a.startHTTP(ctx)
	})

	return grp.Wait()
}

func (a *App) startHTTP(ctx context.Context) error {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Fatal("failed to create listener")
	}

	handler := a.router

	a.httpServer = &http.Server{
		Handler: handler,
	}

	logging.GetLogger(ctx).Info("http server completely initialized and started")

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logging.GetLogger(ctx).Warning("server shutdown")
		default:
			logging.GetLogger(ctx).Fatal(err)
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logging.GetLogger(ctx).Error(err)
	}

	return err
}
