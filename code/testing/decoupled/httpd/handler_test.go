package httpd_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/corylanou/go-presentations/code/testing/decoupled/httpd"
)

func TestSet_NoErrors(t *testing.T) {

	// START MOCK-SETUP-OMIT
	handler := httpd.NewHandler()
	store := &MockStore{}
	store.getFn = func(key string) (interface{}, error) {
		return "bar", nil
	}
	handler.Store = store
	// END MOCK-SETUP-OMIT

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/key", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.Form = url.Values{"key": []string{"foo"}, "value": []string{"bar"}}

	handler.ServeHTTP(w, r)
	if exp, got := w.Code, http.StatusAccepted; exp != got {
		t.Errorf("unexpected error code. exp: %d, got %d", exp, got)
	}
}

// START NO-KEY-OMIT
func TestGet_NoKey(t *testing.T) {
	handler := httpd.NewHandler()
	store := &MockStore{}
	handler.Store = store

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/key", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(w, r)
	if exp, got := w.Code, http.StatusBadRequest; exp != got {
		t.Log(w.Body)
		t.Errorf("unexpected error code. exp: %d, got %d", exp, got)
	}
}

// END NO-KEY-OMIT

// START NOT-FOUND-OMIT
func TestGet_NotFound(t *testing.T) {
	handler := httpd.NewHandler()
	store := &MockStore{}
	store.getFn = func(string) (interface{}, error) {
		return nil, notFound{}
	}
	handler.Store = store

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/key?key=foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(w, r)
	if got, exp := w.Code, http.StatusNotFound; got != exp {
		t.Log(w.Body)
		t.Errorf("unexpected error code. got: %d, exp %d", got, exp)
	}
}

// END NOT-FOUND-OMIT

// START GET-SERVER-ERROR-OMIT
func TestGet_ServerError(t *testing.T) {
	handler := httpd.NewHandler()
	store := &MockStore{}
	store.getFn = func(string) (interface{}, error) {
		return nil, errors.New("boom")
	}
	handler.Store = store

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/key?key=foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(w, r)
	if got, exp := w.Code, http.StatusInternalServerError; got != exp {
		t.Log(w.Body)
		t.Errorf("unexpected error code. got: %d, exp %d", got, exp)
	}
}

// END GET-SERVER-ERROR-OMIT

// START GET-SUCCESS-OMIT
func TestGet_Success(t *testing.T) {
	handler := httpd.NewHandler()
	store := &MockStore{}
	store.getFn = func(string) (interface{}, error) {
		return "bar", nil
	}
	handler.Store = store

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/key?key=foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(w, r)
	if got, exp := w.Code, http.StatusOK; got != exp {
		t.Fatalf("unexpected status code.  got %d, expected %d", got, exp)
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	if got, exp := data["foo"], "bar"; got != exp {
		t.Fatalf("unexpected value.  got: %v, exp %v", got, exp)
	}
}

// END GET-SUCCESS-OMIT

// START MOCK-OMIT
type MockStore struct {
	setFn func(key string, value interface{})
	getFn func(key string) (interface{}, error)
}

func (ms *MockStore) Set(key string, value interface{}) {
	if ms.setFn != nil {
		ms.setFn(key, value)
	}
}

func (ms *MockStore) Get(key string) (interface{}, error) {
	if ms.getFn != nil {
		return ms.getFn(key)
	}
	return nil, nil
}

// END MOCK-OMIT

// START VERBOSE-OMIT
func TestVerbose(t *testing.T) {
	if testing.Verbose() {
		t.Log("put extra logging here that normally we don't care about")
	} else {
		// silence my normal loggers
		log.SetOutput(ioutil.Discard)
	}
}

// END VERBOSE-OMIT

// not found mock
type notFound struct{}

func (nf notFound) NotFound() {}

func (nf notFound) Error() string {
	return ""
}
