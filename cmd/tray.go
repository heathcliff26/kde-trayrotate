package main

import (
	"log/slog"
	"time"

	"fyne.io/systray"
	resources "github.com/heathcliff26/kde-trayrotate"
	"github.com/heathcliff26/kde-trayrotate/pkg/autostart"
	"github.com/heathcliff26/kde-trayrotate/pkg/kscreendoctor"
)

const (
	Title   = "Auto Rotate"
	Tooltip = "Toggle Auto Rotate"
)

var doctor = kscreendoctor.NewKScreenDoctor()
var exit = make(chan struct{})

func run() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle(Title)
	systray.SetOnTapped(toggleAutoRotate)

	systraySetDefault()

	mAutostart := systray.AddMenuItemCheckbox("Autostart", "Start on login", autostart.IsEnabled())
	go toggleAutostartLoop(mAutostart)

	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		for {
			updateSystray()
			select {
			case <-time.Tick(1 * time.Second):
			case <-exit:
				return
			}
		}
	}()
}

func onExit() {
	close(exit)
}

func toggleAutostartLoop(mAutostart *systray.MenuItem) {
	for {
		select {
		case <-mAutostart.ClickedCh:
		case <-exit:
			return
		}
		if mAutostart.Checked() {
			err := autostart.Disable()
			if err != nil {
				slog.Error("Failed to disable autostart", "error", err)
			}
		} else {
			err := autostart.Enable()
			if err != nil {
				slog.Error("Failed to enable autostart", "error", err)
			}
		}

		if autostart.IsEnabled() {
			mAutostart.Check()
		} else {
			mAutostart.Uncheck()
		}
	}
}

func systraySetDefault() {
	systray.SetIcon(resources.Icon)
	systray.SetTooltip(Tooltip)
}

func updateSystray() {
	status, err := doctor.GetStatus()
	if err != nil {
		slog.Error("Failed to get current screen status", "error", err)
		systraySetDefault()
		return
	}

	state, ok := status["1"]
	if !ok {
		slog.Debug("No screens with auto-rotation support found")
		systraySetDefault()
		return
	}

	updateSystrayFromState(state)
}

func updateSystrayFromState(state string) {
	switch state {
	case kscreendoctor.PolicyNever:
		systray.SetIcon(resources.IconRotationLocked)
	case kscreendoctor.PolicyAlways:
		systray.SetIcon(resources.IconRotationUnlocked)
	default:
		systraySetDefault()
	}
}

func toggleAutoRotate() {
	status, err := doctor.GetStatus()
	if err != nil {
		slog.Error("Failed to get current screen status", "error", err)
		systraySetDefault()
		return
	}

	state, ok := status["1"]
	if !ok {
		slog.Debug("No screens with auto-rotation support found")
		systraySetDefault()
		return
	}

	var newPolicy string
	switch state {
	case kscreendoctor.PolicyNever:
		newPolicy = kscreendoctor.PolicyAlways
	case kscreendoctor.PolicyInTabledMode:
		newPolicy = kscreendoctor.PolicyNever
	case kscreendoctor.PolicyAlways:
		newPolicy = kscreendoctor.PolicyNever
	}
	err = doctor.SetAutoRotatePolicies(status, newPolicy)
	if err != nil {
		slog.Error("Failed to set new auto-rotation policy", "error", err)
		systraySetDefault()
		return
	}

	updateSystrayFromState(newPolicy)
}
