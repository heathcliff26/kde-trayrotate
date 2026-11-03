package icons

import (
	_ "embed"
)

//go:embed rotation.png
var Icon []byte

//go:embed rotation-locked.png
var RotationLocked []byte

//go:embed rotation-unlocked.png
var RotationUnlocked []byte
