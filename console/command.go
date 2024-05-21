package console

import (
	"github.com/spf13/cobra"
	"github.com/yusologia/go-core/console/command"
)

func Commands(cobraCmd *cobra.Command, newCommands []BaseInterface) {
	addCommand(cobraCmd, &command.DeleteLogFileCommand{})

	for _, newCommand := range newCommands {
		addCommand(cobraCmd, newCommand)
	}
}

func addCommand(cmd *cobra.Command, newCmd BaseInterface) {
	newCmd.Command(cmd)
}
