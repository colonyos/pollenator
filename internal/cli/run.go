package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/colonyos/colonies/pkg/client"
	"github.com/colonyos/pollinator/pkg/colonies"
	"github.com/colonyos/pollinator/pkg/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&Follow, "follow", "f", false, "Follow process")
	runCmd.Flags().IntVarP(&Count, "count", "", 100, "Number of messages to fetch")
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

		// Sync all directories
		err = colonies.SyncDir("/src", client, ColonyID, ExecutorPrvKey, proj, true)
		CheckError(err)
		err = colonies.SyncDir("/data", client, ColonyID, ExecutorPrvKey, proj, true)
		CheckError(err)

		snapshotID, err := colonies.CreateSrcSnapshot(client, ColonyID, ExecutorPrvKey, proj)
		CheckError(err)

		log.Debug("Generating function spec")
		funcSpec := colonies.CreateFuncSpec(ColonyID, proj, snapshotID)
		CheckError(err)

		addedProcess, err := client.Submit(funcSpec, ExecutorPrvKey)
		CheckError(err)

		url := DashboardURL + "/process?processid=" + addedProcess.ID
		link := fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a\n", url, url)

		log.WithFields(log.Fields{"ProcessID": addedProcess.ID}).Info("Process submitted")
		log.Info("Follow process at " + link)

		if Follow {
			err = colonies.Follow(client, addedProcess, ExecutorPrvKey, Count)
			CheckError(err)
			err = colonies.SyncDir("/result", client, ColonyID, ExecutorPrvKey, proj, false)
			CheckError(err)
		}
	},
}
