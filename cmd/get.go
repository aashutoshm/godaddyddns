package cmd

import (
	"context"
	"log"
	"time"

	"github.com/sickyoon/daddyddns/ddns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get assigned IP",
	Long:  "get assigned IP",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c := ddns.New("y00ns.com")
		resp, err := c.GetCurrentIP(ctx, viper.GetString("subdomain"))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("response: %s", string(resp))
	},
}

func init() {
	refreshCmd.AddCommand(getCmd)
}
