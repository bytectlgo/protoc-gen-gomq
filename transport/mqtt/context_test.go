package mqtt

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	xhttp "github.com/go-kratos/kratos/v2/transport/http"
)

var testRouter = &Router{srv: NewServer(nil)}

func TestContextHeader(t *testing.T) {
	w := wrapper{
		router: testRouter,
		req:    &http.Request{Header: map[string][]string{"name": {"kratos"}}},
		res:    nil,
		w:      responseWriter{},
	}
	h := w.Header()
	if !reflect.DeepEqual(h, http.Header{"name": {"kratos"}}) {
		t.Errorf("expected %v, got %v", http.Header{"name": {"kratos"}}, h)
	}
}

func TestContextForm(t *testing.T) {
	w := wrapper{
		router: testRouter,
		req:    &http.Request{Header: map[string][]string{"name": {"kratos"}}, Method: http.MethodPost},
		res:    nil,
		w:      responseWriter{},
	}
	form := w.Form()
	if !reflect.DeepEqual(form, url.Values{}) {
		t.Errorf("expected %v, got %v", url.Values{}, form)
	}

	w = wrapper{
		router: testRouter,
		req:    &http.Request{Form: map[string][]string{"name": {"kratos"}}},
		res:    nil,
		w:      responseWriter{},
	}
	form = w.Form()
	if !reflect.DeepEqual(form, url.Values{"name": {"kratos"}}) {
		t.Errorf("expected %v, got %v", url.Values{"name": {"kratos"}}, form)
	}
}

func TestContextQuery(t *testing.T) {
	w := wrapper{
		router: testRouter,
		req:    &http.Request{URL: &url.URL{Scheme: "https", Host: "github.com", Path: "go-kratos/kratos", RawQuery: "page=1"}, Method: http.MethodPost},
		res:    nil,
		w:      responseWriter{},
	}
	q := w.Query()
	if !reflect.DeepEqual(q, url.Values{"page": {"1"}}) {
		t.Errorf("expected %v, got %v", url.Values{"page": {"1"}}, q)
	}
}

func TestContextResponse(t *testing.T) {
	res := httptest.NewRecorder()
	w := wrapper{
		router: &Router{srv: &Server{enc: xhttp.DefaultResponseEncoder}},
		req:    &http.Request{Method: http.MethodPost},
		res:    res,
		w:      responseWriter{200, res},
	}
	if !reflect.DeepEqual(w.Response(), res) {
		t.Errorf("expected %v, got %v", res, w.Response())
	}
	err := w.Returns(map[string]string{}, nil)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	needErr := errors.New("some error")
	err = w.Returns(map[string]string{}, needErr)
	if !errors.Is(err, needErr) {
		t.Errorf("expected %v, got %v", needErr, err)
	}
}

func TestResponseUnwrap(t *testing.T) {
	res := httptest.NewRecorder()
	f := func(rw http.ResponseWriter, _ *http.Request, _ interface{}) error {
		u, ok := rw.(interface {
			Unwrap() http.ResponseWriter
		})
		if !ok {
			return errors.New("can not unwrap")
		}
		w := u.Unwrap()
		if !reflect.DeepEqual(w, res) {
			return errors.New("underlying response writer not equal")
		}
		return nil
	}

	w := wrapper{
		router: &Router{srv: &Server{enc: f}},
		req:    nil,
		res:    res,
		w:      responseWriter{200, res},
	}
	err := w.JSON("test", "ok")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
}

func TestContextResponseReturn(t *testing.T) {
	writer := httptest.NewRecorder()
	w := wrapper{
		router: testRouter,
		req:    nil,
		res:    writer,
		w:      responseWriter{},
	}
	err := w.JSON("/test", "success")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}

}

func TestContextCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &http.Request{Method: http.MethodPost}
	req = req.WithContext(ctx)
	w := wrapper{
		router: testRouter,
		req:    req,
		res:    nil,
		w:      responseWriter{},
	}
	_, ok := w.Deadline()
	if !ok {
		t.Errorf("expected %v, got %v", true, ok)
	}
	done := w.Done()
	if done == nil {
		t.Errorf("expected %v, got %v", true, ok)
	}
	err := w.Err()
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	v := w.Value("test")
	if v != nil {
		t.Errorf("expected %v, got %v", nil, v)
	}

	w = wrapper{
		router: &Router{srv: &Server{enc: xhttp.DefaultResponseEncoder}},
		req:    nil,
		res:    nil,
		w:      responseWriter{},
	}
	_, ok = w.Deadline()
	if ok {
		t.Errorf("expected %v, got %v", false, ok)
	}
	done = w.Done()
	if done != nil {
		t.Errorf("expected not nil, got %v", done)
	}
	err = w.Err()
	if err == nil {
		t.Errorf("expected not %v, got %v", nil, err)
	}
	v = w.Value("test")
	if v != nil {
		t.Errorf("expected %v, got %v", nil, v)
	}
}
