package fixture_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/halimath/fixture"
)

type fixtureMock struct {
	label string
	b     *strings.Builder
}

func (f *fixtureMock) BeforeAll(t *testing.T) error {
	fmt.Fprintf(f.b, "%s::BeforeAll", f.label)
	return nil
}

func (f *fixtureMock) AfterAll(t *testing.T) error {
	fmt.Fprintf(f.b, "%s::AfterAll", f.label)
	return nil
}

func (f *fixtureMock) BeforeEach(t *testing.T) error {
	fmt.Fprintf(f.b, "%s::BeforeEach", f.label)
	return nil
}

func (f *fixtureMock) AfterEach(t *testing.T) error {
	fmt.Fprintf(f.b, "%s::AfterEach", f.label)
	return nil
}

func TestFixture_lifecycle(t *testing.T) {
	f := &fixtureMock{
		b: &strings.Builder{},
	}
	t.Run("test", func(t *testing.T) {
		With(t, f).
			Run("test 1", func(t *testing.T, f *fixtureMock) {
				f.b.WriteString("Test#1")
			}).
			Run("test 2", func(t *testing.T, f *fixtureMock) {
				f.b.WriteString("Test#2")
			})
	})

	want := "::BeforeAll::BeforeEachTest#1::AfterEach::BeforeEachTest#2::AfterEach::AfterAll"
	got := f.b.String()

	if want != got {
		t.Errorf("expected %q but got %q", want, got)
	}
}

func TestFixture_tupleFixture(t *testing.T) {
	var b strings.Builder
	f1 := &fixtureMock{label: "f1", b: &b}
	f2 := &fixtureMock{label: "f2", b: &b}

	t.Run("test", func(t *testing.T) {
		With(t, Tuple(f1, f2)).
			Run("test 1", func(t *testing.T, f *TupleFixture[*fixtureMock, *fixtureMock]) {
				f.One.b.WriteString("f1::Test#1")
				f.Two.b.WriteString("f2::Test#1")
			}).
			Run("test 2", func(t *testing.T, f *TupleFixture[*fixtureMock, *fixtureMock]) {
				f.One.b.WriteString("f1::Test#2")
				f.Two.b.WriteString("f2::Test#2")
			})
	})

	want := "f1::BeforeAllf2::BeforeAllf1::BeforeEachf2::BeforeEachf1::Test#1f2::Test#1f2::AfterEachf1::AfterEachf1::BeforeEachf2::BeforeEachf1::Test#2f2::Test#2f2::AfterEachf1::AfterEachf2::AfterAllf1::AfterAll"
	got := b.String()

	if want != got {
		t.Errorf("\nexpected %q\n but got %q", want, got)
	}
}
