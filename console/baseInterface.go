package console

import "github.com/spf13/cobra"

type BaseInterface interface {
	Command(cmd *cobra.Command)
	Handle()
}
