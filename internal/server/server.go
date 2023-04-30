package server

import (
	"context"
	"net/http"

	"jonathanschmitt.com.br/go-api/internal/config"
	exampleService "jonathanschmitt.com.br/go-api/internal/service_example/service"

	"github.com/gorilla/mux"
)

func NewHTTPServer(
	ctx context.Context,
	config config.ServiceConfiguration,
	srvExampleService exampleService.Service,
) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	ConfigureRoutes(ctx, config, r, srvExampleService)

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")
		rw.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
		rw.Header().Add("Access-Control-Expose-Headers", "Authorization")
		next.ServeHTTP(rw, r)
	})
}
