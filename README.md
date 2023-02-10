# fixture

Re-usable test setups und teardowns for go(lang) `testing` tests.

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]

`fixture` implements a micro-framework ontop of the standard library's 
`testing` package that allow writing of reusable test setup and teardown code.

# Installation

This module uses golang modules and can be installed with

```shell
go get github.com/halimath/fixture@main
```

# Usage

`fixture` defines a very simple (and assumingly _well-known_) lifecycle for
code to execute before, inbetween and after tests. A _fixture_ may hook into
this lifecycle to setup or teardown resources needed by the tests. As multiple
tests may share some amount of these resources `fixture` provides a simple
_test suite_ functionality that plays well with the resource initialization.

The lifecycle is shown in the following picture:

![lifecycle](https://www.plantuml.com/plantuml/png/JKyn3i8m3Dpz2giJ8FKBq2AnCR9L7I9mb8ZKgM9twEznWr3nOj-TVROxKLTqcHA0l2FFhhW9xv59rvam5mqPu70wOjkUyKe-5-fJWXtTt3DK-23HMlHUgLGwUcmcwq4rJJ0ooXALBWrg80Qqs0Q6bMJyjwCajAkSnw_djlZ7sab0_8eUeBDi3tm0 "lifecycle")

A _fixture_ (in terms of this package) is any go value. A fixture may satisfy
a couple of additional interfaces to execute code at the given lifecycle
phases. The interfaces are named after the lifecycle phases. Each interface
contains a single method (named after the interface) that receives the
`*testing.T` and returns an `error` which will abort the test (calling
`t.Fatal`).

## Using a fixture

Using a fixture is done using the `With` function, which starts a new test 
suite. Calling `Run` registers a test to run using this fixture.

```go
With(t, new(myFixture)).
	Run("test 1", func(t *testing.T, f *myFixture) {
		// Test code
	}).
	Run("test 2", func(t *testing.T, f *myFixture) {
		// Test code
	})
```

## Implementing a fixture

To implement a fixture simply create a type to hold all the values your fixture
will provide. You can also add receiver functions to ease interaction with the
fixture. Then, implement the desired hook interfaces. 

Typically, a fixture implements the hook methods via a pointer receiver. This
allows using just `new` to create a fixture. Use either `BeforeAll` or
`BeforeEach` to initialize the code.

The following example uses a fixture to spawn a `httptest.Server` with a simple
handler (in a real world the handle would have been some real production code).
It provides a `sendRequest` method to send a simple request, handle errors by
failing the test and returns the `http.Response`. 

The `TestExample` executes two tests both using the same running server.

```go
// A simple test fixture holding a httptest.Server.
type httpServerFixture struct {
	srv *httptest.Server
}

// BeforeAll hooks into the fixture lifecycle and creates and starts the
// httptest.Server before the first test is executed.
func (f *httpServerFixture) BeforeAll(t *testing.T) error {
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Tracing-Id", "1")
		w.WriteHeader(http.StatusNoContent)
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
			if got != http.StatusNoContent {
				t.Errorf("expected %d but got %d", http.StatusNoContent, got)
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
```

## Fixtures already provided by `fixture`

`fixture` contains some ready to use generic fixtures. All these fixtures
have dependencies only to the standard library and cause no external module to
be required.

### TempDir

Creating and removing a temporary directory for filesystem related tests is 
easy with the `TempDirFixture` and the `TempDir` function.

```go
With(t, TempDir("someprefix")).
	Run("create file", func(t *testing.T, d *TempDirFixture) {
		f, err := os.Create(d.Join("test"))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
	}).
	Run("expect file", func(t *testing.T, d *TempDirFixture) {
		_, err := os.Stat(d.Join("test"))
		if err != nil {
			t.Error(err)
		}
	})
```

### HTTPServerFixture

The `HTTPServerFixture` creates a HTTP server using `httptest.NewServer` which
will be started on `BeforeAll` and closed on `AfterAll`. The server uses a
`http.ServerMux` as its handler and handler functions can be registered at any
stage. The server uses HTTP/2 but no TLS; both can be changed easily.

```go
f := new(HTTPServerFixture)

f.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

With(t, f).
	Run("/", func(t *testing.T, f *HTTPServerFixture) {
		res, err := http.Get(f.URL())
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected 200 but got %d", res.StatusCode)
		}
	})
```

# License

Copyright 2022 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[ci-img-url]: https://github.com/halimath/fixture/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/fixture
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/fixture
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/fixture
[release-img-url]: https://img.shields.io/github/v/release/halimath/fixture.svg
[release-url]: https://github.com/halimath/fixture/releases