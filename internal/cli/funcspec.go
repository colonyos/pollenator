package cli

import (
	"encoding/json"
	"fmt"

	"github.com/colonyos/colonies/pkg/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(funcSpecCmd)
	funcSpecCmd.Flags().BoolVarP(&Follow, "follow", "f", false, "Follow process")
}

var funcSpecCmd = &cobra.Command{
	Use:   "funcspec",
	Short: "Sync and generate a Func spec",
	Long:  "Sync and generate a Func spec",
	Run: func(cmd *cobra.Command, args []string) {
		parseEnv()

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		client := client.CreateColoniesClient(ColoniesServerHost, ColoniesServerPort, ColoniesInsecure, ColoniesSkipTLSVerify)

		funcSpec, _ := SyncAndGenerateFuncSpec(client)

		jsonBytes, err := json.MarshalIndent(funcSpec, "", "    ")
		CheckError(err)

		fmt.Println(string(jsonBytes))
	},
}
