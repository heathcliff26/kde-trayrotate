package autostart

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const XDG_CONFIG_HOME = "XDG_CONFIG_HOME"

func TestIsEnabled(t *testing.T) {
	t.Run("Enabled", func(t *testing.T) {
		require := require.New(t)

		t.Setenv(XDG_CONFIG_HOME, t.TempDir())

		require.False(IsEnabled(), "Should not be enabled")
	})
	t.Run("Disabled", func(t *testing.T) {
		require := require.New(t)

		t.Setenv(XDG_CONFIG_HOME, t.TempDir())

		require.NoError(Enable(), "Should enable autostart")
		require.True(IsEnabled(), "Should be enabled")
	})
}

func TestEnable(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		require := require.New(t)

		t.Setenv(XDG_CONFIG_HOME, t.TempDir())

		require.NoError(Enable(), "Should enable autostart")

		path, err := pathToAutostartFile()
		require.NoError(err, "Should get autostart file path")
		res, err := os.ReadFile(path)
		require.NoError(err, "Should read autostart file")
		require.Equal(desktopFile, string(res), "Autostart file content should match")
	})
	t.Run("NoOverwrite", func(t *testing.T) {
		require := require.New(t)

		t.Setenv(XDG_CONFIG_HOME, t.TempDir())

		path, err := pathToAutostartFile()
		require.NoError(err, "Should get autostart file path")
		require.NoError(os.MkdirAll(filepath.Dir(path), 0755), "Should create autostart directory")
		require.NoError(os.WriteFile(path, []byte("test"), 0644), "Should write test content")

		require.NoError(Enable(), "Should not fail when autostart file already exists")

		res, err := os.ReadFile(path)
		require.NoError(err, "Should read autostart file")
		require.Equal("test", string(res), "Autostart file should be unchanged")
	})
}

func TestDisable(t *testing.T) {
	t.Run("Enabled", func(t *testing.T) {
		require := require.New(t)

		t.Setenv(XDG_CONFIG_HOME, t.TempDir())

		require.NoError(Enable(), "Should enable autostart")
		require.NoError(Disable(), "Should disable autostart")
		require.False(IsEnabled(), "Should not be enabled")
	})
	t.Run("NotEnabled", func(t *testing.T) {
		require := require.New(t)

		t.Setenv(XDG_CONFIG_HOME, t.TempDir())
		require.NoError(Disable(), "Should not fail when not enabled")
	})
}

func TestPathToAutostartFile(t *testing.T) {
	require := require.New(t)

	t.Setenv(XDG_CONFIG_HOME, t.TempDir())

	path, err := pathToAutostartFile()
	require.NoError(err, "Should get autostart file path")
	require.Contains(path, filepath.Join("autostart", packageID+".desktop"), "Autostart file path should contain autostart directory and package ID")
}
