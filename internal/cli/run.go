package cli

import (
	"fmt"

	"github.com/colonyos/colonies/pkg/client"
	"github.com/colonyos/pollinator/pkg/colonies"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		client := client.CreateColoniesClient(ColoniesServerHost, ColoniesServerPort, ColoniesInsecure, ColoniesSkipTLSVerify)

		funcSpec, proj := SyncAndGenerateFuncSpec(client)

		addedProcess, err := client.Submit(funcSpec, PrvKey)
		CheckError(err)

		url := DashboardURL + "/process?processid=" + addedProcess.ID
		link := fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a\n", url, url)

		log.WithFields(log.Fields{"ProcessID": addedProcess.ID}).Info("Process submitted")
		log.Info("Follow process at " + link)

		if Follow {
			err = colonies.Follow(client, addedProcess, PrvKey, Count)
			CheckError(err)
			err = colonies.SyncDir("/result", client, ColonyName, PrvKey, proj, false)
			CheckError(err)
		}
	},
}
