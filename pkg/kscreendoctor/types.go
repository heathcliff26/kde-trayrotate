package kscreendoctor

type ScreenSettings struct {
	Outputs []Output `json:"outputs"`
}

type Output struct {
	AutoRotatePolicy *int `json:"autoRotatePolicy,omitempty"`
}
