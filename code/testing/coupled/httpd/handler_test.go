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

	"github.com/corylanou/go-presentations/code/testing/coupled/httpd"
	"github.com/corylanou/go-presentations/code/testing/coupled/keys"
)

func Test_Upsert_Sleep(t *testing.T) {
	handler := httpd.NewHandler()
	store := keys.NewStore()
	handler.Store = store

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

	// START SLEEP-OMIT
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "/key?key=foo", nil)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(4 * time.Second)

	handler.ServeHTTP(w, r)
	if exp, got := w.Code, http.StatusOK; exp != got {
		t.Log(w.Body)
		t.Errorf("unexpected error code. exp: %d, got %d", exp, got)
	}
	// END SLEEP-OMIT
}

func Test_Upsert_Channels(t *testing.T) {
	handler := httpd.NewHandler()
	store := keys.NewStore()
	handler.Store = store

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
				return true, fmt.Errorf("unexpected value.  exp: %s, got %s", exp, got)
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
