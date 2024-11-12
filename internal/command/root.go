package command

import (
	"fmt"
	"log"
	"os"
	"github.com/spf13/cobra"

)

var cliVersion 			string
var nsxApi 					string
var nsxUsername 		string
var nsxPassword 		string
var outputType 			string
var outputFile	string
var skipVerify 			bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: cliVersion,
	Use:     "cfnru",
	Short:   "cfnru blah",
	Long:    "cfnru blah",
	// Example: fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s", downloadUsage, getProductsUsage, getSubProductsUsage, getVersions, getFiles, getManifestExample),
	Run: func(cmd *cobra.Command, args []string) {
		l := log.New(os.Stderr, "", 0)
		validateCredentials(cmd)
		runReport(l)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func runReport(l Logger) {
	generateReport(nsxApi, nsxUsername, nsxPassword, outputType, outputFile, skipVerify, l)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&nsxApi, "nsxApi", "a", "", "IP of FQDN of the NSX Manager [$NSX_API]")
	rootCmd.PersistentFlags().StringVarP(&nsxUsername, "user", "u", "", "Username used to authenticate [$NSX_USER]")
	rootCmd.PersistentFlags().StringVarP(&nsxPassword, "pass", "p", "", "Password used to authenticate [$NSX_PASS]")
	rootCmd.PersistentFlags().StringVarP(&outputType, "type", "t", "json", "Output file type. [json, xlsx]. (default: xlsx)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "file", "f", "", "Output file name. (default: report.xlsx)")
	rootCmd.PersistentFlags().BoolVarP(&skipVerify, "skipVerify", "k", false, "Skip TLS verification")
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}