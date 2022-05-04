package requests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestRequestsTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer server.Close()

	Options(WithTimeout(100 * time.Millisecond))
	t.Run("process timeout check", func(t *testing.T) {
		_, err := process(context.TODO(), server.URL, http.MethodGet, nil)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return
			}
		}
	})
}

func TestRequestsGet(t *testing.T) {
	s1 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`ok`))
	}))
	defer s1.Close()

	s2 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`ok`))
	}))
	defer s2.Close()

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "requests get",
			args:    args{s1.URL},
			want:    []byte(`ok`),
			wantErr: false,
		},
		{
			name:    "requests get StatusBadRequest",
			args:    args{s2.URL},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(context.TODO(), tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestsPost(t *testing.T) {
	s1 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			t.Errorf("invalid method: %s", req.Method)
			return
		}

		var b = struct{ Name string }{}

		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			t.Error(err)
			return
		}
		defer req.Body.Close()

		rw.Write([]byte(b.Name))
	}))
	defer s1.Close()

	type args struct {
		path string
		body []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "requests post successfully",
			args:    args{s1.URL, []byte(`{"name": "ok"}`)},
			want:    []byte(`ok`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(context.TODO(), tt.args.path, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestRequestsDelete(t *testing.T) {
	s1 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			t.Errorf("invalid method: %s", req.Method)
			return
		}
		rw.Write([]byte(`ok`))
	}))
	defer s1.Close()

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "requests delete",
			args:    args{s1.URL},
			want:    []byte(`ok`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Delete(context.TODO(), tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
