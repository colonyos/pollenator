package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const TimeLayout = "2006-01-02 15:04:05"
const KEYCHAIN_PATH = ".colonies"

var Verbose bool
var ColoniesServerHost string
var ColoniesServerPort int
var ColoniesInsecure bool
var ColoniesSkipTLSVerify bool
var ColoniesUseTLS bool
var ColonyID string
var ExecutorID string
var ExecutorPrvKey string
var ExecutorType string
var Follow bool
var Count int

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
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
