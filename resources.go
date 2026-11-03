package resources

import (
	_ "embed"
)

//go:embed packages/io.github.heathcliff26.kde-trayrotate.desktop
var DesktopFile []byte

//go:embed packages/io.github.heathcliff26.kde-trayrotate.png
var Icon []byte

//go:embed icons/rotation-locked.png
var IconRotationLocked []byte

//go:embed icons/rotation-unlocked.png
var IconRotationUnlocked []byte
