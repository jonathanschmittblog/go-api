package v1

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type ApiRequestService interface {
	Token(token string)
	InsecureSkipVerify(skip bool)
	URL(complementPath string) url.URL
	Debug() ApiRequestService
	EndSlash() ApiRequestService
	SetAPM(enable bool)
	Get(ctx context.Context, path string, response interface{}) (int, error)
	Post(ctx context.Context, path string, body interface{}, response interface{}) (int, error)
	Put(ctx context.Context, path string, body interface{}, response interface{}) (int, error)
	Patch(ctx context.Context, path string, body interface{}, response interface{}) (int, error)
	Delete(ctx context.Context, path string, response interface{}) (int, error)
	DeleteWithBody(ctx context.Context, path string, body interface{}, response interface{}) (int, error)
	BasicAuth(username string, password string)
	AddHeaders(headers map[string]string)
	GetBodyResponse() []byte
}

type request struct {
	Url                url.URL
	token              string
	insecureSkipVerify bool
	debug              bool
	endSlash           bool
	context            context.Context
	ElasticAPM         bool
	userInfo           user
	headers            map[string]string
	bodyResponse       []byte
}

type user struct {
	username, password string
}

func (a *request) Token(token string) {
	strArr := strings.Split(token, " ")
	if len(strArr) > 1 {
		token = strArr[len(strArr)-1]
	}
	a.token = token
	a.userInfo.username = ""
	a.userInfo.password = ""
}

func (a *request) GetBodyResponse() []byte {
	return a.bodyResponse
}

func (a *request) SetAPM(enable bool) {
	a.ElasticAPM = enable
}

func (a *request) InsecureSkipVerify(skip bool) {
	a.insecureSkipVerify = skip
}

func (a *request) URL(complementPath string) url.URL {
	uPath, _ := url.Parse(complementPath)

	u := a.Url
	u.Path = path.Join(u.Path, uPath.Path)
	u.RawQuery = uPath.RawQuery
	return u
}

func (a *request) Debug() ApiRequestService {
	b := *a
	b.debug = true
	return &b
}

func (a *request) EndSlash() ApiRequestService {
	b := *a
	b.endSlash = true
	return &b
}

func (a *request) Get(ctx context.Context, path string, response interface{}) (int, error) {
	a.Url.Path = path
	return a.httpRequest(ctx, a.Url, http.MethodGet, nil, response)
}

func (a *request) Post(ctx context.Context, path string, body interface{}, response interface{}) (int, error) {
	a.Url.Path = path
	return a.httpRequest(ctx, a.Url, http.MethodPost, body, response)
}

func (a *request) Put(ctx context.Context, path string, body interface{}, response interface{}) (int, error) {
	a.Url.Path = path
	return a.httpRequest(ctx, a.Url, http.MethodPut, body, response)
}

func (a *request) Patch(ctx context.Context, path string, body interface{}, response interface{}) (int, error) {
	a.Url.Path = path
	return a.httpRequest(ctx, a.Url, http.MethodPatch, body, response)
}

func (a *request) Delete(ctx context.Context, path string, response interface{}) (int, error) {
	a.Url.Path = path
	return a.httpRequest(ctx, a.Url, http.MethodDelete, nil, response)
}

func (a *request) DeleteWithBody(ctx context.Context, path string, body interface{}, response interface{}) (int, error) {
	a.Url.Path = path
	return a.httpRequest(ctx, a.Url, http.MethodDelete, body, response)
}

func (a *request) BasicAuth(username string, password string) {
	a.userInfo = user{
		username,
		password,
	}
	a.token = ""
}

func (a *request) AddHeaders(headers map[string]string) {
	a.headers = headers
}

func (a *request) httpRequest(ctx context.Context, u url.URL, method string, body interface{}, response interface{}) (status int, err error) {
	a.bodyResponse = []byte{}
	do := func() ([]byte, int, error) {
		urlStr := u.Scheme + "://" + u.Host + u.Path
		if a.endSlash {
			urlStr += "/"
		}

		if u.RawQuery != "" {
			urlStr += "?" + u.RawQuery
		}

		var payload string
		if body != nil {
			data, err := json.Marshal(body)
			if err != nil {
				return nil, 0, err
			}
			payload = string(data)
		}

		client := &http.Client{}
		if a.insecureSkipVerify {
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		}

		req, err := http.NewRequest(method, urlStr, bytes.NewBufferString(payload))
		if err != nil {
			return nil, 0, err
		}

		for key, value := range a.headers {
			req.Header.Set(key, value)
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		if (len(a.userInfo.username) > 0) && (len(a.userInfo.password) > 0) {
			req.SetBasicAuth(a.userInfo.username, a.userInfo.password)
		} else if len(a.token) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
		}

		if http.MethodPut == method || http.MethodPost == method || http.MethodPatch == method || (http.MethodDelete == method && body != nil) {
			req.Header.Set("Content-Type", "application/json")
			req.ParseForm()
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, 0, err
		}
		defer resp.Body.Close()

		bodyResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, err
		}

		if a.debug && !(http.MethodGet == method && (strings.Contains(req.URL.Path, "/jwt") || strings.Contains(req.URL.Path, "/users"))) {
			fmt.Printf("\n%s | %s %s | Body: %s", req.Method, resp.Status, req.URL, string(bodyResponse))
		}

		if resp.StatusCode >= 400 {
			return bodyResponse, resp.StatusCode, fmt.Errorf(string(bodyResponse))
		}

		if resp.StatusCode == 204 {
			return bodyResponse, resp.StatusCode, nil
		}

		return bodyResponse, resp.StatusCode, nil
	}

	bodyResponse, statusCode, err := do()
	if err != nil {
		return statusCode, err
	}
	a.bodyResponse = bodyResponse

	if response != nil && bodyResponse != nil {
		if json.Valid(bodyResponse) {
			err = json.Unmarshal(bodyResponse, response)
			if err != nil {
				return statusCode, err
			}
		} else {
			err = yaml.Unmarshal(bodyResponse, response)
			if err != nil {
				return statusCode, err
			}
		}
	}

	return 0, nil
}

func NewAPIRequest(url url.URL) ApiRequestService {
	ctx := context.TODO()

	request := request{
		Url:     url,
		context: ctx,
	}

	return &request
}
