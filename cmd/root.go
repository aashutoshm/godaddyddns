package cmd

import (
	"fmt"
	"log"
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
	cobra.OnInitialize()
	viper.SetEnvPrefix("ddns")
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringP("access_key", "a", viper.GetString("access_key"), "GoDaddy API Access Key")
	viper.BindPFlag("access_key", rootCmd.PersistentFlags().Lookup("access_key"))
	rootCmd.PersistentFlags().StringP("secret_key", "s", viper.GetString("secret_key"), "GoDaddy API Secret Key")
	viper.BindPFlag("secret_key", rootCmd.PersistentFlags().Lookup("secret_key"))

	rootCmd.PersistentFlags().BoolP("ote", "e", false, "GoDaddy Environment (Production by default)")
	viper.BindPFlag("ote", rootCmd.PersistentFlags().Lookup("ote"))

	// validate required args
	if viper.GetString("access_key") == "" || viper.GetString("secret_key") == "" {
		log.Fatalf("access_key or secret_key must not be empty")
	}
}

// Execute command line
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
