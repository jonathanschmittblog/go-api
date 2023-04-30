package server

import (
	"context"
	"fmt"
	"net/http"

	core "jonathanschmitt.com.br/go-api/pkg/core/v1"

	config "jonathanschmitt.com.br/go-api/internal/config"

	exampleService "jonathanschmitt.com.br/go-api/internal/service_example/service"

	"github.com/gorilla/mux"
)

func ConfigureRoutes(ctx context.Context, config config.ServiceConfiguration, r *mux.Router,
	srvExample exampleService.Service,
) {
	sr := r.PathPrefix(fmt.Sprintf("/api/%s", config.Server.Version)).Subrouter()
	defaultRoutes(sr)
	exampleServiceRoutes(ctx, config, sr, srvExample)
}

func defaultRoutes(r *mux.Router) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		core.ResponseSuccess(w, "heath")
	}).Methods("GET")
}

func exampleServiceRoutes(ctx context.Context, config config.ServiceConfiguration, r *mux.Router, srvExample exampleService.Service) {
	r.HandleFunc("/service-example/{id}", srvExample.Get).Methods("GET")
	r.HandleFunc("/service-example", srvExample.List).Methods("GET")
	r.HandleFunc("/service-example", srvExample.Add).Methods("POST")
	r.HandleFunc("/service-example/{id}", srvExample.Update).Methods("PUT")
	r.HandleFunc("/service-example/{id}", srvExample.Delete).Methods("DELETE")
}
