package main

import (
	"testing"
)

func TestIsSevenZipInstalled(tester *testing.T) {

	tester.Run("doesNotError", func(test *testing.T) {
		_, err := isSevenZipInstalled()
		if err != nil {
			test.Errorf("got %q, want nil", err.Error())
		}
	})
}
