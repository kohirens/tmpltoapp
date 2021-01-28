package stdlib

import (
	"os"
	"runtime"
	"testing"
)

func TestHomeDir(t *testing.T) {
	nonWindowsWant, windowsWant := "/home/"+os.Getenv("USER"), "C:\\Users\\"+os.Getenv("USERNAME")

	t.Run("success", func(t *testing.T) {
		// exec code.
		got, gotErr := HomeDir()
		// the want changes from Windows to *nix systems.
		if runtime.GOOS == "windows" {
			if got != windowsWant {
				t.Errorf("got %q, want %q", got, windowsWant)
			}
		} else if gotErr != nil {
			t.Errorf("not on windows and got %q, but want %q", gotErr.Error(), nonWindowsWant)
		}
	})
}

func TestAppDataDir(t *testing.T) {
	nonWindowsWant, windowsWant := "/home/"+os.Getenv("USER"), "C:\\Users\\"+os.Getenv("USERNAME")+"\\AppData\\Local"

	t.Run("success", func(t *testing.T) {
		got, err := AppDataDir()

		if err != nil {
			t.Errorf("unexpeted error during test, got %q", err.Error())
		}

		if runtime.GOOS == "windows" {
			if got != windowsWant {
				t.Errorf("want %q, got %q", windowsWant, got)
			}
		} else { // Linux/Mac
			if got != nonWindowsWant {
				t.Errorf("want %q, got %q", nonWindowsWant, got)
			}
		}
	})
}
