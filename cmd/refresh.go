package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "refresh DNS record",
	Long:  "refresh DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		// TOOD:
		// PATCH https://api.godaddy.com/v1/domains/{domain}/records
		/*
			[
				{
					"data": "string",
					"name": "string",
					"port": 0,
					"priority": 0,
					"protocol": "string",
					"service": "string",
					"ttl": 0,
					"type": "A",
					"weight": 0
				}
			]
		*/
	},
}

func init() {
	refreshCmd.PersistentFlags().StringP("subdomain", "d", "dev", "Subdomain")
	viper.BindPFlag("subdomain", refreshCmd.PersistentFlags().Lookup("subdomain"))
	rootCmd.AddCommand(refreshCmd)
}
