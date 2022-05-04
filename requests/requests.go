package requests

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/palantir/stacktrace"
)

// var client = &http.Client{Timeout: 5 * time.Second}

type client struct {
	*http.Client
}

type opt func(*http.Client) error
type header func(http.Header)

func NewClient() (*client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}
	return &client{
		&http.Client{
			Timeout: 5 * time.Second,
			Jar:     jar,
		},
	}, nil
}

func WithTimeout(x time.Duration) opt {
	return func(c *http.Client) error {
		c.Timeout = x
		return nil
	}
}

func WithCookieJar() opt {
	return func(c *http.Client) error {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}
		c.Jar = jar
		return nil
	}
}

func WithHeader(key, value string) header {
	return func(h http.Header) {
		h.Set(key, value)
	}
}

func (c *client) do(ctx context.Context, urlPath string, method string, values url.Values, h ...header) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, urlPath, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	for _, v := range h {
		v(req.Header)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, stacktrace.NewError("%s %s", resp.Status, data)
	}

	return data, nil
}

func (c *client) Get(ctx context.Context, path string, headers ...header) ([]byte, error) {
	return c.do(ctx, path, http.MethodGet, nil, headers...)
}

func (c *client) Post(ctx context.Context, path string, values url.Values, headers ...header) ([]byte, error) {
	return c.do(ctx, path, http.MethodPost, values, headers...)
}

func (c *client) Delete(ctx context.Context, path string, headers ...header) ([]byte, error) {
	return c.do(ctx, path, http.MethodDelete, nil, headers...)
}
