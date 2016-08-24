package httpd_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/corylanou/go-presentations/code/testing/decoupled/httpd"
)

func Test_Upsert_NoErrors(t *testing.T) {

	// START MOCK-SETUP-OMIT
	handler := httpd.NewHandler()
	store := &mockStore{}
	store.getFn = func(key string) (interface{}, error) {
		return "bar", nil
	}
	handler.Store = store
	// END MOCK-SETUP-OMIT

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/key", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.Form = url.Values{"key": []string{"foo"}, "value": []string{"bar"}}

	handler.ServeHTTP(w, r)
	if exp, got := w.Code, http.StatusAccepted; exp != got {
		t.Errorf("unexpected error code. exp: %d, got %d", exp, got)
	}

	// START FUNC-INNER-OMIT
	test := func() (bool, error) {
		w = httptest.NewRecorder()
		r, err = http.NewRequest("GET", "/key?key=foo", nil)
		if err != nil {
			return false, err
		}
		handler.ServeHTTP(w, r)

		if exp, got := w.Code, http.StatusOK; exp == got {
			data := map[string]interface{}{}
			if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
				return false, err
			}
			if exp, got := "bar", data["foo"]; exp != got {
				return true, fmt.Errorf("unexpected value.  exp: %v, got %v", exp, got)
			} else {
				return true, nil
			}
		}
		return false, nil
	}
	// END FUNC-INNER-OMIT

	// START CHANNEL-SETUP-OMIT
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.NewTimer(4 * time.Second)
	defer timeout.Stop()
	// END CHANNEL-SETUP-OMIT

	// START CHANNEL-OMIT
	for {
		select {
		case <-timeout.C:
			t.Fatal("test timed out waiting for success")
			return
		case <-ticker.C:
			ok, err := test()
			if ok && err != nil {
				t.Fatal(err)
				return
			}
			if ok && err == nil {
				// test successful
				return
			}
		}
	}
	// END CHANNEL-OMIT
}

// START MOCK-OMIT
type mockStore struct {
	upsertFn func(key string, value interface{})
	getFn    func(key string) (interface{}, error)
}

func (ms *mockStore) Upsert(key string, value interface{}) {
	if ms.upsertFn != nil {
		ms.upsertFn(key, value)
	}
}

func (ms *mockStore) Get(key string) (interface{}, error) {
	if ms.getFn != nil {
		return ms.getFn(key)
	}
	return nil, nil
}

// END MOCK-OMIT
