package internal

import (
	"database/sql"

	"jonathanschmitt.com.br/go-api/internal/config"

	exampleRepositories "jonathanschmitt.com.br/go-api/internal/service_example/repositories"

	"github.com/go-kit/log"
)

type Repositories struct {
	ExampleRepo exampleRepositories.Repository
}

func InitializeRepositories(db *sql.DB, logger log.Logger, config config.ServiceConfiguration) Repositories {
	return Repositories{
		ExampleRepo: exampleRepositories.NewRepo(db, logger, config),
	}
}
