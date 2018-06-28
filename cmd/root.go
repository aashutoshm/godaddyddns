package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ddns",
	Short: "GoDaddy DDNS",
	Long:  "GoDaddy DDNS - dynamic dns through GoDaddy API",
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("access_key", "a", viper.GetString("access_key"), "GoDaddy API Access Key")
	viper.BindPFlag("access_key", rootCmd.PersistentFlags().Lookup("access_key"))
	rootCmd.PersistentFlags().StringP("secret_key", "s", viper.GetString("secret_key"), "GoDaddy API Secret Key")
	viper.BindPFlag("secret_key", rootCmd.PersistentFlags().Lookup("secret_key"))
}

func initConfig() {
	viper.SetEnvPrefix("ddns")
	viper.AutomaticEnv()
}

// Execute command line
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
