package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/colonyos/colonies/pkg/client"
	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/fs"
	"github.com/colonyos/pollinator/pkg/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Deploy and run project",
	Long:  "Deploy and run project",
	Run: func(cmd *cobra.Command, args []string) {
		parseEnv()

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		projectFile := "./project.yaml"
		projectData, err := ioutil.ReadFile(projectFile)
		CheckError(err)

		proj := &project.Project{}
		err = yaml.Unmarshal([]byte(projectData), &proj)
		CheckError(err)

		client := client.CreateColoniesClient(ColoniesServerHost, ColoniesServerPort, ColoniesInsecure, ColoniesSkipTLSVerify)

		log.Debug("Starting a file storage client")
		fsClient, err := fs.CreateFSClient(client, ColonyID, ExecutorPrvKey)
		CheckError(err)

		keepLocal := true
		label := "/pollinator/" + proj.ProjectID + "/src"
		syncPlans, err := fsClient.CalcSyncPlans("./cfs/src", label, keepLocal)
		CheckError(err)

		fmt.Println(syncPlans)

		// Create snapshot
		snapshotID := core.GenerateRandomID()

		log.WithFields(log.Fields{
			"SnapshotID": snapshotID,
			"Dir":        "./cfs/src"}).
			Info("Creating snapshot")

		log.Info("Generating function spec")
		funcSpec := project.CreateFuncSpec(ColonyID, proj, snapshotID)
		CheckError(err)
		jsonString, err := funcSpec.ToJSON()
		CheckError(err)
		fmt.Println(jsonString)
	},
}
