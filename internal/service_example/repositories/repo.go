package service_example

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	config "jonathanschmitt.com.br/go-api/internal/config"
	core "jonathanschmitt.com.br/go-api/pkg/core/v1"
	model "jonathanschmitt.com.br/go-api/pkg/service_example/v1"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Repository interface {
	Get(ctx context.Context, id string) (model.Example, error)
	List(ctx context.Context, name string, limit int, actualPage int) (core.Paging, error)
	Add(ctx context.Context, example *model.Example, tx *sql.Tx) error
	Update(ctx context.Context, example *model.Example, tx *sql.Tx) error
	Delete(ctx context.Context, example *model.Example, tx *sql.Tx) error
	GetDB() *sql.DB
}

type repo struct {
	db     *sql.DB
	logger log.Logger
	config config.ServiceConfiguration
}

func NewRepo(db *sql.DB, logger log.Logger, config config.ServiceConfiguration) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
		config: config,
	}
}

func (repo *repo) GetDB() *sql.DB {
	return repo.db
}

func (repo *repo) Get(ctx context.Context, id string) (model.Example, error) {
	result := model.Example{Id: id}

	sqlString := `
SELECT et.id, et.name
FROM example_table et
WHERE et.id = $1`

	var name sql.NullString
	err := repo.db.QueryRow(sqlString, id).Scan(&result.Id, &name)
	if err != nil {
		return result, err
	}
	if name.Valid {
		result.Name = name.String
	}

	return result, nil
}

func (repo *repo) List(ctx context.Context, name string, limit int, actualPage int) (core.Paging, error) {
	result := []model.Example{}
	var params []interface{}

	sqlWhere := ""
	if len(name) > 0 {
		sqlWhere = sqlWhere + "AND et.name = $" + strconv.Itoa(len(params)+1) + " "
		params = append(params, name)
	}
	if sqlWhere != "" {
		sqlWhere = "WHERE" + sqlWhere[3:]
	}

	sqlString := `
	SELECT 
		et.id, et.name, count(*) OVER() AS full_count
	FROM example_table et	
	<##WHERE##>
	ORDER  BY
		et.name DESC
	OFFSET <##OFFSET##>
	LIMIT  <##LIMIT##>`

	if limit <= 0 {
		limit = repo.config.Server.PageSize
	}
	if actualPage <= 0 {
		actualPage = 1
	}

	sqlString = strings.Replace(sqlString, "<##WHERE##>", sqlWhere, 1)
	sqlString = strings.Replace(sqlString, "<##OFFSET##>", strconv.Itoa(limit*(actualPage-1)), 1)
	sqlString = strings.Replace(sqlString, "<##LIMIT##>", strconv.Itoa(limit), 1)

	sqlStatement, err := repo.db.Query(sqlString, params...)
	if err != nil {
		return core.NewPaging(0, 0, 0, result), err
	}
	defer sqlStatement.Close()

	var full_count int
	for sqlStatement.Next() {
		item := model.Example{}
		err = sqlStatement.Scan(&item.Id, &item.Name, &full_count)
		if err != nil {
			return core.NewPaging(0, 0, 0, result), err
		}
		result = append(result, item)
	}

	return core.NewPaging(actualPage, full_count, limit, result), nil
}

func (repo *repo) Add(ctx context.Context, example *model.Example, tx *sql.Tx) error {
	logger := log.With(repo.logger, "method", "Add")

	sql := `INSERT INTO example_table (name) VALUES ($1) RETURNING id`
	err := tx.QueryRow(sql, example.Name).Scan(&example.Id)
	if err != nil {
		level.Error(logger).Log(err)
		return err
	}
	logger.Log("new example id: ", example.Id)

	return nil
}

func (repo *repo) Update(ctx context.Context, example *model.Example, tx *sql.Tx) error {
	sqlString := `
	UPDATE example_table
	SET name = $1
	WHERE id = $2`

	_, err := tx.Exec(sqlString, example.Name, example.Id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repo) Delete(ctx context.Context, example *model.Example, tx *sql.Tx) error {
	sql := `
	DELETE FROM example_table
	WHERE id = $1`

	_, err := tx.Exec(sql, example.Id)
	if err != nil {
		return err
	}

	return nil
}
