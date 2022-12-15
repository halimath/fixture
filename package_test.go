package fixture_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/halimath/fixture"
)

// A simple test fixture holding a httptest.Server.
type httpServerFixture struct {
	srv *httptest.Server
}

// BeforeAll hooks into the fixture lifecycle and creates and starts the
// httptest.Server before the first test is executed.
func (f *httpServerFixture) BeforeAll(t *testing.T) error {
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Tracing-Id", "1")
		w.WriteHeader(http.StatusFound)
	}))

	return nil
}

// AfterAll hooks into the fixture lifecycle and disposes the httptest.Server
// after the last test has been executed.
func (f *httpServerFixture) AfterAll(t *testing.T) error {
	f.srv.Close()
	return nil
}

// sendRequest is a convenience function making it easier to read the test code.
func (f *httpServerFixture) sendRequest(t *testing.T) *http.Response {
	r, err := http.Get(f.srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestExample(t *testing.T) {
	fixture.With(t, new(httpServerFixture)).
		Run("http status code", func(t *testing.T, f *httpServerFixture) {
			got := f.sendRequest(t).StatusCode
			if got != http.StatusFound {
				t.Errorf("expected %d but got %d", http.StatusFound, got)
			}
		}).
		Run("tracing header", func(t *testing.T, f *httpServerFixture) {
			resp := f.sendRequest(t)
			got := resp.Header.Get("X-Tracing-Id")
			if got != "1" {
				t.Errorf("expected %q but got %q", "1", got)
			}
		})
}
