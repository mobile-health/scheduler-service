package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/mobile-health/scheduler-service/src/models"
	"goji.io"
)

type Client struct {
	BaseApiURL string
	Mux        *goji.Mux
	ApiKey     string
	ApiLogin   string
}

func (c *Client) perform(r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c.Mux.ServeHTTP(w, r)
	r.Body.Close()
	return w
}

func (c *Client) do(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	req.SetBasicAuth(c.ApiKey, c.ApiLogin)

	return c.perform(req)
}

func (c *Client) DoPost(endpoint string, body io.Reader) (models.MapInterface, error) {
	apiURL := c.BaseApiURL + endpoint

	resp := c.do(http.MethodPost, apiURL, body)
	if resp.Code != 201 {
		return nil, errors.New(resp.Body.String())
	}

	var result models.MapInterface

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
