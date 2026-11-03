package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/heathcliff26/kde-trayrotate/pkg/version"
	"github.com/spf13/cobra"
)

const flagLogLevel = "log"

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
	}

	logLevel := rootCmd.Flags().StringP(flagLogLevel, "l", "info", "Set the log level")

	rootCmd.AddCommand(
		version.NewCommand(),
	)

	rootCmd.Run = func(cmd *cobra.Command, _ []string) {
		err := setLogLevel(*logLevel)
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			os.Exit(1)
		}
		run()
	}

	return rootCmd
}

// Parse a given string and set the resulting log level
func setLogLevel(level string) error {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		return fmt.Errorf("invalid log level '%s'", level)
	}

	opts := slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)

	return nil
}
