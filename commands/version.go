package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Gets current newsfetch version",
	Run: func(cmd *cobra.Command, args []string) {
		if VERSION == "" {
			fmt.Println("Could not find version number, package must be built with `make build`")
		} else {
			fmt.Println("Version: ", VERSION)
		}

		if COMMITHASH == "" {
			fmt.Println("Could not find git commit hash, package must be built with `make build`")
		} else {
			fmt.Println("Git commit: ", COMMITHASH)
		}
	},
}
