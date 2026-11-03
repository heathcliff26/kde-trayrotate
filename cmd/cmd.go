package main

import (
	"strings"

	"github.com/heathcliff26/kde-trayrotate/pkg/version"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	name := strings.ToLower(version.Name)
	cobra.AddTemplateFunc(
		"ProgramName", func() string {
			return name
		},
	)

	rootCmd := &cobra.Command{
		Use:   name,
		Short: version.Name + " toggle auto rotate from the system tray",
		Run: func(_ *cobra.Command, _ []string) {
			run()
		},
	}

	rootCmd.AddCommand(
		version.NewCommand(),
	)

	return rootCmd
}
