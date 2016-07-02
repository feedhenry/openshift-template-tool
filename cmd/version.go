package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	version = "dev build"
	commit  = "dev build"
)

func VersionCommand(stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersion(stdout, stderr)
		},
		SilenceUsage: true,
	}
	return cmd
}

func RunVersion(stdout, stderr io.Writer) error {
	fmt.Fprintf(stdout, "Version: %v \nCommit: %v\n", version, commit)
	return nil
}
