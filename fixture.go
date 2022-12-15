// Package fixture provides a micro-framework on top of the testing package
// that provides test setup and teardown handling in a reusable way. The
// package defines a (well-known) lifecycle to execute, which is in detail
// documented in the repo's README file.
//
// A fixture is any go value. To hook into the lifecycle a fixture must
// satisfy any of the single-method interfaces defined.
package fixture

import (
	"sync"
	"testing"
)

// Fixture defines a "label" interface for fixtures providing values for tests.
// This interface defines no methods and is completely equivalent to any. It
// just serves the purpose of making the generic type annotations easier to
// reason about.
type Fixture interface{}

// BeforeAll is an extension hook interface that defines the BeforeAll hook.
type BeforeAll interface {
	Fixture
	BeforeAll(t *testing.T) error
}

// AfterAll is an extension hook interface that defines the AfterAll hook.
type AfterAll interface {
	Fixture
	AfterAll(t *testing.T) error
}

// BeforeEach is an extension hook interface that defines the BeforeEach hook.
type BeforeEach interface {
	BeforeEach(t *testing.T) error
}

// AfterEach is an extension hook interface that defines the AfterEach hook.
type AfterEach interface {
	AfterEach(t *testing.T) error
}

// TestFunc defines the type for test functions running a test on behalf of
// a Fixture.
type TestFunc[F Fixture] func(*testing.T, F)

// Suite defines a suite of tests using the same fixture.
type Suite[F Fixture] interface {
	Run(string, TestFunc[F]) Suite[F]
}

// suiteBuilder is an implementation of a Suite.
type suiteBuilder[F Fixture] struct {
	f              F
	t              *testing.T
	testRun        bool
	afterAllRunner sync.Once
}

// Run registers another test with name to run using the fixture contained in f.
// It makes sure that all hooks are executed according to the lifecycle and
// executes test on behalf of f.t.
func (f *suiteBuilder[F]) Run(name string, test TestFunc[F]) Suite[F] {
	f.t.Helper()

	var fix any = f.f

	if aa, ok := fix.(AfterAll); ok {
		f.afterAllRunner.Do(func() {
			f.t.Helper()
			f.t.Cleanup(func() {
				f.t.Helper()

				if err := aa.AfterAll(f.t); err != nil {
					f.t.Fatal(err)
				}
			})
		})
	}

	if !f.testRun {
		if ba, ok := fix.(BeforeAll); ok {
			if err := ba.BeforeAll(f.t); err != nil {
				f.t.Fatal(err)
			}
		}
	}
	f.testRun = true

	if ba, ok := fix.(BeforeEach); ok {
		if err := ba.BeforeEach(f.t); err != nil {
			f.t.Fatal(err)
		}
	}

	f.t.Run(name, func(t *testing.T) {
		test(t, f.f)
	})

	if ba, ok := fix.(AfterEach); ok {
		if err := ba.AfterEach(f.t); err != nil {
			f.t.Fatal(err)
		}
	}

	return f
}

// With is used to define a new Suite based on fixture.
func With[F any](t *testing.T, fixture F) Suite[F] {
	return &suiteBuilder[F]{
		f: fixture,
		t: t,
	}
}

// TupleFixture is a convenience type used to combine exactly two fixtures into
// a single one to use with tests. It uses MultiFixture under the hood but
// wraps it inside a struct that exposes the fixtures statically typed.
type TupleFixture[A, B Fixture] struct {
	One A
	Two B
	m   MultiFixture
}

func Tuple[A, B Fixture](a A, b B) *TupleFixture[A, B] {
	return &TupleFixture[A, B]{One: a, Two: b}
}

func (f *TupleFixture[A, B]) BeforeAll(t *testing.T) error {
	f.m = MultiFixture{f.One, f.Two}
	return f.m.BeforeAll(t)
}

func (f *TupleFixture[A, B]) AfterAll(t *testing.T) error {
	return f.m.AfterAll(t)
}

func (f *TupleFixture[A, B]) BeforeEach(t *testing.T) error {
	return f.m.BeforeEach(t)
}

func (f *TupleFixture[A, B]) AfterEach(t *testing.T) error {
	return f.m.AfterEach(t)
}

// MultiFixture combines multiple fixtures into a single one to use with With.
// It implements every hook interface and delegates all hooks to each fixture.
// The order of delegation is defined by the hooks type:
//
//   - Before-hooks are execute first-to-last order
//   - After-hooks are executed in last-to-first order
type MultiFixture []Fixture

func (f MultiFixture) BeforeAll(t *testing.T) error {
	for i := range f {
		if h, ok := f[i].(BeforeAll); ok {
			if err := h.BeforeAll(t); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f MultiFixture) AfterAll(t *testing.T) error {
	for i := len(f) - 1; i >= 0; i-- {
		if h, ok := f[i].(AfterAll); ok {
			if err := h.AfterAll(t); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f MultiFixture) BeforeEach(t *testing.T) error {
	for i := range f {
		if h, ok := f[i].(BeforeEach); ok {
			if err := h.BeforeEach(t); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f MultiFixture) AfterEach(t *testing.T) error {
	for i := len(f) - 1; i >= 0; i-- {
		if h, ok := f[i].(AfterEach); ok {
			if err := h.AfterEach(t); err != nil {
				return err
			}
		}
	}
	return nil
}
