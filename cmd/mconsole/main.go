package main

import (
	"fmt"
	"strings"

	"github.com/b177y/go-uml-utilities/pkg/mconsole"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:                   "uml_mconsole socket-path [command]",
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			err := mconsole.RunShell(args[0])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			command := strings.Join(args[1:], " ")
			output, err := mconsole.CommandWithSock(command, args[0])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("OK: %s\n", output)
			}
		}
	},
}

func main() {
	cmd.Execute()
}
