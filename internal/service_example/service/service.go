package service_example

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"jonathanschmitt.com.br/go-api/internal"
	"jonathanschmitt.com.br/go-api/internal/config"
	repositories "jonathanschmitt.com.br/go-api/internal/service_example/repositories"
	core "jonathanschmitt.com.br/go-api/pkg/core/v1"
	model "jonathanschmitt.com.br/go-api/pkg/service_example/v1"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

type Service interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type service struct {
	repository repositories.Repository
	logger     log.Logger
	config     config.ServiceConfiguration
}

func NewService(rep internal.Repositories, logger log.Logger, config config.ServiceConfiguration) Service {
	return &service{
		repository: rep.ExampleRepo,
		logger:     logger,
		config:     config,
	}
}

func (s service) Get(w http.ResponseWriter, r *http.Request) {
	logger := log.With(s.logger, "service", "example", "method", "Get")

	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) <= 0 {
		core.ResponseError(w, core.Response_serviceExampleName, core.Response_idRequired, nil, logger)
		return
	}

	example, err := s.repository.Get(r.Context(), id)
	if err != nil {
		core.ResponseNoContent(w, core.Response_serviceExampleName, core.Response_idNotFound, err, logger)
		return
	}
	core.ResponseSuccess(w, example)
}

func (s service) List(w http.ResponseWriter, r *http.Request) {
	logger := log.With(s.logger, "service", "example", "method", "Get")

	actualPage, err := strconv.Atoi(r.FormValue("actualPage"))
	if err != nil || actualPage < 1 {
		actualPage = 1
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil || limit < 1 {
		limit = 0
	}

	listExamples, err := s.repository.List(r.Context(), r.FormValue("name"), limit, actualPage)
	if err != nil {
		core.ResponseDefaultError(w, core.GenericError_Database_Query, err, logger)
		return
	}
	core.ResponseSuccess(w, listExamples)
}

func (s service) Update(w http.ResponseWriter, r *http.Request) {
	logger := log.With(s.logger, "service", "example", "method", "Update")

	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) <= 0 {
		core.ResponseError(w, core.Response_serviceExampleName, core.Response_idRequired, nil, logger)
		return
	}

	example, err := s.repository.Get(r.Context(), id)
	if err != nil {
		core.ResponseNoContent(w, core.Response_serviceExampleName, core.Response_idNotFound, err, logger)
		return
	}

	var newExample model.Example
	reqBody, err := io.ReadAll(r.Body)
	if err != nil || !json.Valid(reqBody) {
		core.ResponseDefaultError(w, core.JsonNotGetFromBody, err, logger)
		return
	}
	if err = json.Unmarshal(reqBody, &newExample); err != nil {
		core.ResponseDefaultError(w, core.JsonNotDeserialized, err, logger)
		return
	}

	if len(newExample.Name) <= 0 {
		core.ResponseError(w, core.Response_serviceExampleName, core.Response_nameRequired, nil, logger)
		return
	}

	example.Name = newExample.Name

	tx, err := s.repository.GetDB().BeginTx(context.Background(), nil)
	if err != nil {
		core.ResponseDefaultError(w, core.GenericError_Unexpected, err, logger)
		return
	}
	defer tx.Rollback()

	if err = s.repository.Update(r.Context(), &example, tx); err != nil {
		core.ResponseDefaultError(w, core.GenericError_Database_Persistence, err, logger)
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		core.ResponseDefaultError(w, core.GenericError_Unexpected, err, logger)
		return
	}

	logger.Log(fmt.Sprintf("example id: '%s' updated", example.Id))
	core.ResponseSimpleSuccess(w)
}

func (s service) Add(w http.ResponseWriter, r *http.Request) {
	logger := log.With(s.logger, "service", "example", "method", "Add")

	var example model.Example
	reqBody, err := io.ReadAll(r.Body)
	if err != nil || !json.Valid(reqBody) {
		core.ResponseDefaultError(w, core.JsonNotGetFromBody, err, logger)
		return
	}
	if err = json.Unmarshal(reqBody, &example); err != nil {
		core.ResponseDefaultError(w, core.JsonNotDeserialized, err, logger)
		return
	}

	if len(example.Name) <= 0 {
		core.ResponseError(w, core.Response_serviceExampleName, core.Response_nameRequired, nil, logger)
		return
	}

	tx, err := s.repository.GetDB().BeginTx(context.Background(), nil)
	if err != nil {
		core.ResponseDefaultError(w, core.GenericError_Unexpected, err, logger)
		return
	}
	defer tx.Rollback()

	err = s.repository.Add(r.Context(), &example, tx)
	if err != nil {
		core.ResponseDefaultError(w, core.GenericError_Database_Persistence, err, logger)
		return
	}

	if err = tx.Commit(); err != nil {
		core.ResponseDefaultError(w, core.GenericError_Unexpected, err, logger)
		return
	}
	core.ResponseSuccess(w, example)
}

func (s service) Delete(w http.ResponseWriter, r *http.Request) {
	logger := log.With(s.logger, "service", "example", "method", "Delete")

	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) <= 0 {
		core.ResponseError(w, core.Response_serviceExampleName, core.Response_idRequired, nil, logger)
		return
	}

	example, err := s.repository.Get(r.Context(), id)
	if err != nil {
		core.ResponseNoContent(w, core.Response_serviceExampleName, core.Response_idNotFound, err, logger)
		return
	}

	tx, err := s.repository.GetDB().BeginTx(context.Background(), nil)
	if err != nil {
		core.ResponseDefaultError(w, core.GenericError_Unexpected, err, logger)
		return
	}
	defer tx.Rollback()

	if err = s.repository.Delete(r.Context(), &example, tx); err != nil {
		core.ResponseDefaultError(w, core.GenericError_Database_Persistence, err, logger)
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		core.ResponseDefaultError(w, core.GenericError_Unexpected, err, logger)
		return
	}

	logger.Log(fmt.Sprintf("example id: '%s' removed", example.Id))
	core.ResponseSimpleSuccess(w)
}
