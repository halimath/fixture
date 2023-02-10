package fixture

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

// HTTPServerFixture is a fixture that provides a httptest.Server for testing.
// The server will be started on BeforeAll and closed on AfterAll. The server
// is started with HTTP2 enabled but without TLS by default. Both can be
// changed by setting the boolean flags on the fixture.
type HTTPServerFixture struct {
	mux          *http.ServeMux
	srv          *httptest.Server
	UseTLS       bool
	DisableHTTP2 bool
}

func (f *HTTPServerFixture) Handle(pattern string, handler http.Handler) {
	if f.mux == nil {
		f.mux = http.NewServeMux()
	}
	f.mux.Handle(pattern, handler)
}

func (f *HTTPServerFixture) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if f.mux == nil {
		f.mux = http.NewServeMux()
	}
	f.mux.HandleFunc(pattern, handler)

}

// URL returns the url used to connect to the test server formed by using the
// server's base url (which contains the randomly chosen port as well as the
// respective protocol) and all elements from pathElements joined with a /.
func (f *HTTPServerFixture) URL(pathElements ...string) string {
	if len(pathElements) == 0 {
		return f.srv.URL + "/"
	}

	p := path.Join(pathElements...)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	return f.srv.URL + p
}

func (f *HTTPServerFixture) BeforeAll(t *testing.T) error {
	f.srv = httptest.NewUnstartedServer(f.mux)
	f.srv.EnableHTTP2 = !f.DisableHTTP2

	if f.UseTLS {
		f.srv.StartTLS()
	} else {
		f.srv.Start()
	}

	return nil
}

func (f *HTTPServerFixture) AfterAll(t *testing.T) error {
	f.srv.Close()
	return nil
}
