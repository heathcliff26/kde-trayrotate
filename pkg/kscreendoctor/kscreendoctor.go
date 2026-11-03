package kscreendoctor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
)

const (
	autoRotatePolicyTemplate = "output.%s.autoRotatePolicy.%s"

	PolicyNever        = "never"
	PolicyInTabledMode = "inTabletMode"
	PolicyAlways       = "always"
)

type KScreenDoctor struct {
	// Binary is only configurable for testing purposes.
	binary string
}

// Create a new instance for kscreen-doctor commands
func NewKScreenDoctor() *KScreenDoctor {
	return &KScreenDoctor{
		binary: "kscreen-doctor",
	}
}

// Get the current status of the screen rotation settings
func (ksc *KScreenDoctor) GetStatus() (map[string]string, error) {
	out, err := ksc.runCommand("--json")
	if err != nil {
		return nil, fmt.Errorf("failed to get settings from kscreen-doctor: %w", err)
	}

	var settings ScreenSettings
	err = json.Unmarshal([]byte(out), &settings)
	if err != nil {
		return nil, fmt.Errorf("failed to parse settings from kscreen-doctor: %w", err)
	}

	if len(settings.Outputs) == 0 {
		return nil, fmt.Errorf("no outputs found")
	}

	status := make(map[string]string)

	for i, output := range settings.Outputs {
		if output.AutoRotatePolicy != nil {
			policy := ""
			switch *output.AutoRotatePolicy {
			case 0:
				policy = PolicyNever
			case 1:
				policy = PolicyInTabledMode
			case 2:
				policy = PolicyAlways
			default:
				slog.Debug("Unknown auto-rotation state", slog.Int("output", i), slog.Int("state", *output.AutoRotatePolicy))
			}
			// The screen index used by kscreen-doctor counts from 1 -> i+1
			status[strconv.Itoa(i+1)] = policy
		}
	}

	return status, nil
}

// Set the auto-rotation policy for a given output
func (ksc *KScreenDoctor) SetAutoRotatePolicies(outputs map[string]string, policy string) error {
	for output, state := range outputs {
		if state == policy {
			continue
		}

		_, err := ksc.runCommand(fmt.Sprintf(autoRotatePolicyTemplate, output, policy))
		if err != nil {
			return fmt.Errorf("failed to set auto-rotation policy for output %s: %w", output, err)
		}
	}
	return nil
}

func (ksc *KScreenDoctor) runCommand(args ...string) (string, error) {
	cmd := exec.Command(ksc.binary, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
