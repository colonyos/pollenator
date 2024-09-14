package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const TimeLayout = "2006-01-02 15:04:05"
const KEYCHAIN_PATH = ".colonies"

var ASCII bool
var Verbose bool
var ColoniesServerHost string
var ColoniesServerPort int
var ColoniesInsecure bool
var ColoniesSkipTLSVerify bool
var ColoniesUseTLS bool
var ColonyName string
var PrvKey string
var ExecutorName string
var Follow bool
var Count int
var DashboardURL string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(versionCmd)
}

var rootCmd = &cobra.Command{
	Use:   "pollinator",
	Short: "ColonyOS Pollinator",
	Long:  "ColonyOS Pollinator",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version",
	Long:  "Version",
	Run: func(cmd *cobra.Command, args []string) {
		ASCII = false
		ASCIIStr := os.Getenv("POLLINATOR_CLI_ASCII")
		if ASCIIStr == "true" {
			ASCII = true
		}

		printVersionTable()
	},
}
