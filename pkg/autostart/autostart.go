package autostart

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/heathcliff26/kde-trayrotate/pkg/version"
)

const packageID = "io.github.heathcliff26." + version.Name

// Should be kept in sync with the main .desktop file.
var desktopFile = `[Desktop Entry]
Type=Application

Name=kde-trayrotate
Comment=Quicksetting for auto rotation in the system tray
Categories=Settings;
Keywords=auto-rotate;

Icon=io.github.heathcliff26.kde-trayrotate
Exec=kde-trayrotate
Terminal=false
`

func IsEnabled() bool {
	path, err := pathToAutostartFile()
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return !os.IsNotExist(err)
}

// Create an autostart entry if it doesn't exist.
func Enable() error {
	if IsEnabled() {
		return nil
	}

	path, err := pathToAutostartFile()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(desktopFile), 0644)
}

// Remove the autostart entry if it exists.
func Disable() error {
	path, err := pathToAutostartFile()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func pathToAutostartFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "autostart", packageID+".desktop"), nil
}
