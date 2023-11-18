package cli

import (
	"os"

	"github.com/colonyos/pollinator/pkg/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&ExecutorType, "executortype", "e", "", "Executor type")
	newCmd.MarkFlagRequired("executortype")
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new project",
	Long:  "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		CheckError(checkIfDirIsEmpty("."))
		CheckError(checkIfDirExists("src"))
		CheckError(checkIfDirExists("data"))
		CheckError(checkIfDirExists("result"))

		log.WithFields(log.Fields{
			"Dir": "./src"}).
			Info("Creating directory")
		err := os.MkdirAll("./src", 0755)
		CheckError(err)

		log.WithFields(log.Fields{
			"Dir": "./cfs/data"}).
			Info("Creating directory")
		err = os.MkdirAll("./cfs/data", 0755)

		CheckError(err)

		log.WithFields(log.Fields{
			"Dir": "./cfs/results"}).
			Info("Creating directory")
		err = os.MkdirAll("./cfs/result", 0755)

		CheckError(err)

		err = project.GenerateProjectConfig(ExecutorType)
		CheckError(err)

		err = project.GenerateProjectData()
		CheckError(err)
	},
}
