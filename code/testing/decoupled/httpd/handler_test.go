package httpd_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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

	// START FUNC-INNER-OMIT
	test := func() error {
		w = httptest.NewRecorder()
		r, err = http.NewRequest("GET", "/key?key=foo", nil)
		if err != nil {
			return err
		}
		handler.ServeHTTP(w, r)

		if got, exp := w.Code, http.StatusOK; got != exp {
			return fmt.Errorf("unexpected status code.  got %d, expected %d", got, exp)
		}
		data := map[string]interface{}{}
		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			return err
		}
		if got, exp := data["foo"], "bar"; got != exp {
			return fmt.Errorf("unexpected value.  got: %v, exp %v", got, exp)
		}
		// test successful
		return nil
	}
	// END FUNC-INNER-OMIT

	// START CHANNEL-SETUP-OMIT
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.NewTimer(4 * time.Second)
	defer timeout.Stop()
	// END CHANNEL-SETUP-OMIT

	// START CHANNEL-OMIT
	var testErr error

	for {
		select {
		case <-timeout.C:
			t.Fatalf("test timed out waiting for success.  last error: %s", testErr)
			return
		case <-ticker.C:
			testErr = test()
			if testErr == nil {
				// test successful
				return
			}
		}
	}
	// END CHANNEL-OMIT
}

// START MOCK-OMIT
type MockStore struct {
	upsertFn func(key string, value interface{})
	getFn    func(key string) (interface{}, error)
}

func (ms *MockStore) Set(key string, value interface{}) {
	if ms.upsertFn != nil {
		ms.upsertFn(key, value)
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
