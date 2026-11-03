package main

import (
	"fmt"
	"log/slog"
	"time"

	"fyne.io/systray"
	"github.com/heathcliff26/kde-trayrotate/icons"
	"github.com/heathcliff26/kde-trayrotate/pkg/kscreendoctor"
)

const (
	Title   = "io.github.heathcliff26.kde-trayrotate"
	Tooltip = "Auto Rotate"

	AutoRotateOn  = "on"
	AutoRotateOff = "off"
)

var doctor = kscreendoctor.NewKScreenDoctor()
var exit = make(chan struct{})

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle(Title)
	systray.SetOnTapped(toggleAutoRotate)

	systraySetDefault()

	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	// Sets the icon of a menu item.
	quitCh := mQuit.ClickedCh
	go func() {
		<-quitCh
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

func systraySetDefault() {
	systray.SetIcon(icons.Icon)
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
		systray.SetIcon(icons.RotationLocked)
		systray.SetTooltip(fmt.Sprintf("%s: %s", Tooltip, AutoRotateOff))
	case kscreendoctor.PolicyAlways:
		systray.SetIcon(icons.RotationUnlocked)
		systray.SetTooltip(fmt.Sprintf("%s: %s", Tooltip, AutoRotateOn))
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
