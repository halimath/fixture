package fixture

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type TempDirFixture struct {
	name, path string
}

func TempDir(name string) *TempDirFixture {
	return &TempDirFixture{
		name: name,
	}
}

func (f *TempDirFixture) BeforeAll(t *testing.T) error {
	path, err := os.MkdirTemp("", f.name)
	if err != nil {
		return fmt.Errorf("failed to create temporary test directory %s: %w", f.name, err)
	}
	f.path = path
	return nil
}

func (f *TempDirFixture) AfterAll(t *testing.T) error {
	if err := os.RemoveAll(f.path); err != nil {
		return fmt.Errorf("failed to remove temporary test directory %s: %w", f.path, err)
	}
	return nil
}

func (f *TempDirFixture) Name() string {
	return f.name
}

func (f *TempDirFixture) Path() string {
	return f.path
}

func (f *TempDirFixture) String() string {
	return f.path
}

func (f *TempDirFixture) Join(parts ...string) string {
	p := make([]string, len(parts)+1)
	p[0] = f.path
	copy(p[1:], parts)
	return filepath.Join(p...)
}
