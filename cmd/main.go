package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"jonathanschmitt.com.br/go-api/internal"
	configs "jonathanschmitt.com.br/go-api/internal/config"
	"jonathanschmitt.com.br/go-api/internal/server"
	core "jonathanschmitt.com.br/go-api/pkg/core/v1"

	_ "github.com/lib/pq"

	exampleService "jonathanschmitt.com.br/go-api/internal/service_example/service"

	"github.com/go-kit/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/go-kit/log/level"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "Example api",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	config, err := configs.GetConfig()
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}
	var httpAddr = flag.String("http", fmt.Sprintf(":%d", config.Server.Port), "http listen address")

	level.Info(logger).Log("msg", "service starting")
	defer level.Info(logger).Log("msg", "service ended")

	db, err := core.ConnectDatabase(config)
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	m, err := migrate.NewWithDatabaseInstance("file://../database/migrations/", "postgres", driver)
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	flag.Parse()
	ctx := context.Background()

	respositories := internal.InitializeRepositories(db, logger, config)
	srvExample := exampleService.NewService(respositories, logger, config)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		handler := server.NewHTTPServer(
			ctx,
			config,
			srvExample,
		)
		level.Info(logger).Log("msg", "service started")
		fmt.Println("👷🏼 app listening on port 🚧", *httpAddr, " 🚧")
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
