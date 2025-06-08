package command

import (
	"github.com/spf13/cobra"
	logiapkg "github.com/yusologia/go-core/v2/pkg"
	"os"
	"strconv"
	"time"
)

type DeleteLogFileCommand struct{}

func (command *DeleteLogFileCommand) Command(cmd *cobra.Command) {
	cmd.AddCommand(&cobra.Command{
		Use:  "delete-log-file",
		Long: "Delete log file command",
		Run: func(cmd *cobra.Command, args []string) {
			logiapkg.InitDevMode()

			command.Handle()
		},
	})
}

func (command *DeleteLogFileCommand) Prepare() (cancel func()) {
	return func() {}
}

func (command *DeleteLogFileCommand) Handle() {
	storageDir := os.Getenv("STORAGE_DIR") + "/logs/"

	logDays := 14
	logDaysEnv := os.Getenv("LOG_DAYS")
	if len(logDaysEnv) > 0 {
		logDays, _ = strconv.Atoi(logDaysEnv)
	}

	filename := time.Now().AddDate(0, 0, -logDays).Format("2006-01-02") + ".log"
	fullPath := storageDir + filename
	logiapkg.LogDebug(fullPath)

	_, err := os.Stat(fullPath)
	if err == nil {
		err := os.Remove(fullPath)
		if err != nil {
			logiapkg.LogError(err)
		}
	}
}
