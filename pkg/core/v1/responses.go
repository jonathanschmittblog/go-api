package v1

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/moshenahmias/failure"
)

type errorResponse struct {
	status  int
	message string
	fields  string
}

func ResponseDefaultError(w http.ResponseWriter, key string, err error, logger log.Logger) {
	ResponseError(w, "default", string(key), err, logger)
}

func ResponseError(w http.ResponseWriter, service string, key string, err error, logger log.Logger) {
	if logger != nil && err != nil {
		level.Error(logger).Log("err", err)
	}
	w.WriteHeader(http.StatusBadRequest)
	response := getErrorResponse(getKey(service, key))
	responseMap := map[string]interface{}{}
	responseMap["status"] = response.status
	if err == nil {
		responseMap["message"] = response.message
	} else {
		responseMap["error"] = failure.Build(response.message).ParentOf(err).Done()
	}
	if len(response.fields) > 0 {
		responseMap["fields"] = response.fields
	}
	json.NewEncoder(w).Encode(responseMap)
}

func ResponseNoContent(w http.ResponseWriter, service string, responseKey string, err error, logger log.Logger) {
	if logger != nil && err != nil {
		level.Error(logger).Log("err", err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func ResponseSuccess(w http.ResponseWriter, response interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ResponseSimpleSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func getErrorResponse(key string) errorResponse {
	if errResp, ok := errorsMap.Errors[key]; ok {
		return errResp
	}
	return errorResponse{status: -999, message: "errorNotMapped"}
}

type mapErrorResponse struct {
	sync.RWMutex
	Errors map[string]errorResponse
}

const (
	Response_serviceExampleName = "response_service_example"

	GenericError_Database_Persistence = "GenericError_Database_Persistence"
	GenericError_Database_Query       = "GenericError_Database_Query"
	GenericError_Unexpected           = "GenericError_Unexpected"
	JsonNotGetFromBody                = "JsonNotGetFromBody"
	JsonNotDeserialized               = "JsonNotDeserialized"
	IdNotReceivedFromQuerystring      = "IdNotReceivedFromQuerystring"
	Response_idRequired               = "response_idRequired"
	Response_nameRequired             = "response_nameRequired"
	Response_idNotFound               = "response_idNotFound"
)

var errorsMap = mapErrorResponse{Errors: make(map[string]errorResponse)}

func RegisterErrorResponse(service string, key string, status int, message string, fields string) {
	errorsMap.Lock()
	defer errorsMap.Unlock()
	errorsMap.Errors[getKey(service, key)] = errorResponse{status: status, message: message, fields: fields}
}

func getKey(service string, key string) string {
	return service + "_" + key
}

func init() {
	RegisterErrorResponse("default", string(GenericError_Database_Persistence), 1, "DatabasePersitenceError", "")
	RegisterErrorResponse("default", string(GenericError_Database_Query), 2, "DatabaseQueryError", "")
	RegisterErrorResponse("default", string(GenericError_Unexpected), 3, "UnexpectedError", "")
	RegisterErrorResponse("default", string(JsonNotGetFromBody), 4, "JsonNotGetFromBody", "")
	RegisterErrorResponse("default", string(JsonNotDeserialized), 5, "JsonNotDeserialized", "")
	RegisterErrorResponse("default", string(IdNotReceivedFromQuerystring), 6, "IdNotReceivedFromQuerystring", "")

	//example service
	RegisterErrorResponse(Response_serviceExampleName, Response_idRequired, 7, "id is required", "id")
	RegisterErrorResponse(Response_serviceExampleName, Response_nameRequired, 8, "name is required", "name")
	RegisterErrorResponse(Response_serviceExampleName, Response_idNotFound, 9, "id not found", "id")
}
