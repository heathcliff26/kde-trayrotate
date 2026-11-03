package kscreendoctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewKScreenDoctor(t *testing.T) {
	require := require.New(t)

	doctor := NewKScreenDoctor()
	require.NotNil(doctor, "Should return instance")
	require.Equal("kscreen-doctor", doctor.binary, "Should have correct binary")
}

func TestGetStatus(t *testing.T) {
	tMatrix := []struct {
		Name   string
		Data   string
		Fail   bool
		Error  string
		Result map[string]string
	}{
		{
			Name:   "NoAutoRotate",
			Data:   "testdata/no-auto-rotate.json",
			Result: map[string]string{},
		},
		{
			Name: "AutoRotate",
			Data: "testdata/4-screens.json",
			Result: map[string]string{
				"1": PolicyNever,
				"2": PolicyInTabledMode,
				"3": PolicyAlways,
			},
		},
		{
			Name:  "ExecFail",
			Fail:  true,
			Error: "failed to get settings from kscreen-doctor:",
		},
		{
			Name:  "InvalidJSON",
			Data:  "testdata/invalid-json.txt",
			Error: "failed to parse settings from kscreen-doctor:",
		},
		{
			Name:  "NoOutputs",
			Data:  "testdata/empty.json",
			Error: "no outputs found",
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			require := require.New(t)

			doctor := NewKScreenDoctor()
			doctor.binary = "testdata/kscreen-doctor.sh"

			logfile := filepath.Join(t.TempDir(), "out.txt")
			t.Setenv("KDE_TRAYROTATE_DOCTOR_LOG", logfile)
			t.Setenv("KDE_TRAYROTATE_DOCTOR_DATA", tCase.Data)
			if tCase.Fail {
				t.Setenv("KDE_TRAYROTATE_DOCTOR_FAIL", "true")
			}

			status, err := doctor.GetStatus()

			cmds, cmdErr := os.ReadFile(logfile)
			require.NoError(cmdErr, "Should read logfile")
			require.Equal("--json\n", string(cmds), "Should only call --json")

			if tCase.Error != "" {
				require.ErrorContains(err, tCase.Error, "Should return correct error")
				return
			}

			require.NoError(err, "Should not return error")
			require.Equal(tCase.Result, status, "Should return correct status")
		})
	}
}

func TestSetAutoRotatePolicies(t *testing.T) {
	t.Run("NoOutputs", func(t *testing.T) {
		require := require.New(t)

		doctor := NewKScreenDoctor()
		doctor.binary = "testdata/kscreen-doctor.sh"

		logfile := filepath.Join(t.TempDir(), "out.txt")
		t.Setenv("KDE_TRAYROTATE_DOCTOR_LOG", logfile)

		err := doctor.SetAutoRotatePolicies(nil, PolicyNever)
		require.NoError(err, "Should not fail")

		_, err = os.Stat(logfile)
		require.ErrorIs(err, os.ErrNotExist, "Should not call kscreen-doctor")
	})
	t.Run("SkipAlreadySet", func(t *testing.T) {
		require := require.New(t)

		doctor := NewKScreenDoctor()
		doctor.binary = "testdata/kscreen-doctor.sh"

		logfile := filepath.Join(t.TempDir(), "out.txt")
		t.Setenv("KDE_TRAYROTATE_DOCTOR_LOG", logfile)

		err := doctor.SetAutoRotatePolicies(map[string]string{
			"1": PolicyAlways,
			"2": PolicyNever,
			"3": PolicyInTabledMode,
		},
			PolicyNever,
		)
		require.NoError(err, "Should not fail")

		buf, err := os.ReadFile(logfile)
		require.NoError(err, "Should read logfile")
		cmds := strings.Split(string(buf), "\n")

		if cmds[len(cmds)-1] == "" {
			cmds = cmds[:len(cmds)-1]
		}

		expectedCommands := []string{
			fmt.Sprintf(autoRotatePolicyTemplate, "1", PolicyNever),
			fmt.Sprintf(autoRotatePolicyTemplate, "3", PolicyNever),
		}

		require.ElementsMatch(expectedCommands, cmds, "Should call correct commands")
	})
	t.Run("ExecFail", func(t *testing.T) {
		require := require.New(t)

		doctor := NewKScreenDoctor()
		doctor.binary = "testdata/kscreen-doctor.sh"

		logfile := filepath.Join(t.TempDir(), "out.txt")
		t.Setenv("KDE_TRAYROTATE_DOCTOR_LOG", logfile)
		t.Setenv("KDE_TRAYROTATE_DOCTOR_FAIL", "true")

		err := doctor.SetAutoRotatePolicies(map[string]string{"1": PolicyNever}, PolicyAlways)
		require.ErrorContains(err, "failed to set auto-rotation policy for output 1:", "Should return correct error")

		cmds, err := os.ReadFile(logfile)
		require.NoError(err, "Should read logfile")
		require.Equal(fmt.Sprintf(autoRotatePolicyTemplate, "1", PolicyAlways)+"\n", string(cmds), "Should call correct command")
	})
}
