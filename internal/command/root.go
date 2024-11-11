package command

import (
	"github.com/spf13/cobra"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/collect_data"
)

var nsxApi string
var nsxUsername string
var nsxPassword string
var cliVersion string
var skipVerify bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: cliVersion,
	Use:     "cfnru",
	Short:   "cfnru blah",
	Long:    "cfnru blah",
	// Example: fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s", downloadUsage, getProductsUsage, getSubProductsUsage, getVersions, getFiles, getManifestExample),
	Run: func(cmd *cobra.Command, args []string) {
		validateCredentials(cmd)
		run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func run() {
	collect_data.Run(nsxApi, nsxUsername, nsxPassword, skipVerify)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&nsxApi, "nsxApi", "a", "", "IP of FQDN of the NSX Manager [$NSX_API]")
	rootCmd.PersistentFlags().StringVarP(&nsxUsername, "user", "u", "", "Username used to authenticate [$NSX_USER]")
	rootCmd.PersistentFlags().StringVarP(&nsxPassword, "pass", "p", "", "Password used to authenticate [$NSX_PASS]")
	rootCmd.PersistentFlags().BoolVarP(&skipVerify, "skipVerify", "k", false, "Skip TLS verification")
}