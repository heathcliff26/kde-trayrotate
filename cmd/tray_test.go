package main

import (
	"testing"
	"testing/synctest"

	"fyne.io/systray"
	"github.com/heathcliff26/kde-trayrotate/pkg/autostart"
	"github.com/stretchr/testify/assert"
)

func TestToggleAutostartLoop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		assert := assert.New(t)
		exit = make(chan struct{})

		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		mItem := systray.AddMenuItemCheckbox("Test", "Test", false)

		go toggleAutostartLoop(mItem)
		t.Cleanup(func() {
			close(exit)
			exit = make(chan struct{})
		})

		assert.False(autostart.IsEnabled(), "Should not be enabled at the start")

		mItem.ClickedCh <- struct{}{}
		synctest.Wait()
		assert.True(autostart.IsEnabled(), "Should have enabled autostart")
		assert.True(mItem.Checked(), "Should have ticked checkbox")

		mItem.ClickedCh <- struct{}{}
		synctest.Wait()
		assert.False(autostart.IsEnabled(), "Should have disabled autostart")
		assert.False(mItem.Checked(), "Should have unticked checkbox")
	})
}
