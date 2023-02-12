package cli

import (
	"os"
	"testing"
)

const (
	FixtureDir = "testdata"
	TmpDir     = "tmp"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	_ = os.RemoveAll(TmpDir)
	// Set up a temporary dir for generate files
	_ = os.Mkdir(TmpDir, DirMode) // set up a temporary dir for generate files
	// Run all tests
	exitCode := m.Run()
	// Clean up
	os.Exit(exitCode)
}
