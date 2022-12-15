package fixture

import (
	"errors"
	"os"
	"testing"
)

func TestTempDir(t *testing.T) {
	d := TempDir("temp_dir_test")
	t.Run("test", func(t *testing.T) {
		With(t, d).
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
	})

	_, err := os.Stat(d.Join("test"))
	t.Log(d.Join("test"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Error(err)
	}
}
