package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	core "jonathanschmitt.com.br/go-api/pkg/core/v1"
)

type Service interface {
	Get(r *http.Request, ctx context.Context, id string) (Example, error)
	List(r *http.Request, ctx context.Context, page, limit int) (core.Paging, error)
	Add(r *http.Request, ctx context.Context, example Example) (Example, error)
	Update(r *http.Request, ctx context.Context, example Example) (SimpleSuccessResponse, error)
	Delete(r *http.Request, ctx context.Context, id string) (SimpleSuccessResponse, error)
}

type service struct {
	api core.ApiRequestService
}

func NewService(url_auth_service string) (Service, error) {
	url, err := url.Parse(url_auth_service)
	if err != nil {
		return nil, err
	}
	return &service{
		api: core.NewAPIRequest(*url).Debug(),
	}, nil
}

func (s service) Get(r *http.Request, ctx context.Context, id string) (Example, error) {
	s.api.Token(r.Header.Get("Authorization"))

	var result Example
	url := fmt.Sprintf("/api/v1/service-example/%s", id)
	_, err := s.api.Get(ctx, url, &result)
	if err != nil {
		return Example{}, err
	}
	return result, nil
}

func (s service) List(r *http.Request, ctx context.Context, page, limit int) (core.Paging, error) {
	s.api.Token(r.Header.Get("Authorization"))

	var result core.Paging
	url := fmt.Sprintf("/api/v1/service-example?page=%d&limit=%d", page, limit)
	_, err := s.api.Get(ctx, url, &result)
	if err != nil {
		return core.Paging{}, err
	}
	return result, nil
}

func (s service) Add(r *http.Request, ctx context.Context, example Example) (Example, error) {
	s.api.Token(r.Header.Get("Authorization"))

	url := "/api/v1/service-example"
	var result Example
	_, err := s.api.Post(ctx, url, example, &result)
	if err != nil {
		return Example{}, err
	}
	return result, nil
}

func (s service) Update(r *http.Request, ctx context.Context, example Example) (SimpleSuccessResponse, error) {
	s.api.Token(r.Header.Get("Authorization"))

	vars := r.URL.Query()
	id, _ := strconv.Atoi(vars.Get("id"))
	vars.Del("id")

	url := fmt.Sprintf("/api/v1/service-example/%d", id)
	var result SimpleSuccessResponse
	_, err := s.api.Put(ctx, url, example, &result)
	if err != nil {
		return SimpleSuccessResponse{}, err
	}
	return result, nil
}

func (s service) Delete(r *http.Request, ctx context.Context, id string) (SimpleSuccessResponse, error) {
	s.api.Token(r.Header.Get("Authorization"))

	url := fmt.Sprintf("/api/v1/service-example/%s", id)
	var result SimpleSuccessResponse
	_, err := s.api.Delete(ctx, url, &result)
	if err != nil {
		return SimpleSuccessResponse{}, err
	}
	return result, nil
}
