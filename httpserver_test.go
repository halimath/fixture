package fixture

import (
	"net/http"
	"testing"
)

func TestHTTPServerFixture(t *testing.T) {
	f := new(HTTPServerFixture)

	f.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	f.Handle("/foo/bar", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	With(t, f).
		Run("/", func(t *testing.T, f *HTTPServerFixture) {
			res, err := http.Get(f.URL())
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected 200 but got %d", res.StatusCode)
			}
		}).
		Run("/foo/bar", func(t *testing.T, f *HTTPServerFixture) {
			res, err := http.Get(f.URL("foo", "bar"))
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != http.StatusNoContent {
				t.Errorf("expected 204 but got %d", res.StatusCode)
			}
		})
}
