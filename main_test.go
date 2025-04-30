package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/B1NARY-GR0UP/nwa/cmd"
)

func TestAddLicenseHeader(t *testing.T) {
	temp, cleanup := tempDir(t, "testdata/temp")
	defer cleanup()

	initialDir := "testdata/add/initial"
	expectedDir := "testdata/add/expected"

	if _, err := os.Stat(initialDir); os.IsNotExist(err) {
		t.Skip("Test data directory does not exist")
	}

	err := cp(temp, initialDir)
	if err != nil {
		t.Fatalf("Failed to copy files to temporary directory: %v", err)
	}

	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"nwa", "add", "-c", "BINARY Members", "-y", "2025", temp + "/**"}

	cmd.Execute()

	compareResults(t, temp, expectedDir)
}

func TestCheckLicenseHeader(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"nwa", "check", "-c", "BINARY Members", "-y", "2025", "testdata/check/**"}

	// nwa will exit 1 if check failed
	cmd.Execute()
}

func TestUpdateLicenseHeader(t *testing.T) {
	temp, cleanup := tempDir(t, "testdata/temp")
	defer cleanup()

	initialDir := "testdata/update/initial"
	expectedDir := "testdata/update/expected"

	if _, err := os.Stat(initialDir); os.IsNotExist(err) {
		t.Skip("Test data directory does not exist")
	}

	err := cp(temp, initialDir)
	if err != nil {
		t.Fatalf("Failed to copy files to temporary directory: %v", err)
	}

	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"nwa", "update", "-c", "BINARY Members", "-y", "2025", temp + "/**"}

	cmd.Execute()

	compareResults(t, temp, expectedDir)
}

func TestRemoveLicenseHeader(t *testing.T) {
	temp, cleanup := tempDir(t, "testdata/temp")
	defer cleanup()

	initialDir := "testdata/remove/initial"
	expectedDir := "testdata/remove/expected"

	if _, err := os.Stat(initialDir); os.IsNotExist(err) {
		t.Skip("Test data directory does not exist")
	}

	err := cp(temp, initialDir)
	if err != nil {
		t.Fatalf("Failed to copy files to temporary directory: %v", err)
	}

	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"nwa", "remove", "-c", "BINARY Members", "-y", "2025", temp + "/**"}

	cmd.Execute()

	compareResults(t, temp, expectedDir)
}

func tempDir(t *testing.T, dir string) (string, func()) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Logf("Failed to clean up temporary directory: %v", err)
		}
	}

	return dir, cleanup
}

func compareResults(t *testing.T, actualDir, expectedDir string) {
	err := filepath.Walk(expectedDir, func(expectedPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(expectedDir, expectedPath)
		if err != nil {
			return err
		}

		actualPath := filepath.Join(actualDir, relPath)

		expectedContent, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Errorf("Failed to read expected file %s: %v", expectedPath, err)
			return nil
		}

		actualContent, err := os.ReadFile(actualPath)
		if err != nil {
			t.Errorf("Failed to read processed file %s: %v", actualPath, err)
			return nil
		}

		expectedNormalized := normalizeLine(string(expectedContent))
		actualNormalized := normalizeLine(string(actualContent))

		if expectedNormalized != actualNormalized {
			t.Errorf("File %s content doesn't match\nExpected:\n%s\nActual:\n%s",
				relPath, string(expectedContent), string(actualContent))
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to compare files: %v", err)
	}
}

// normalizeLine CR and CRLF line separator to LF
func normalizeLine(s string) string {
	// crlf -> lf
	s = strings.ReplaceAll(s, "\r\n", "\n")
	// cr -> lf
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func cp(dst, src string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}
