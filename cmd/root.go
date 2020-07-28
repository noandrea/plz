package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool
var welcome string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "plz",
	Short: "",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	rootCmd.Version = v
	// setup the banner
	welcome = fmt.Sprintf(`
╔═╗╦  ╔═╗┌─┐┌─┐┬
╠═╝║  ╔═╝├─┤├─┘│
╩  ╩═╝╚═╝┴ ┴┴  ┴ version %s
	`, rootCmd.Version)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// for debug logging
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if debug {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
}
