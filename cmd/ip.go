package cmd

import (
	"context"
	"log"
	"time"

	"github.com/sickyoon/daddyddns/ddns"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "get external IP",
	Long:  "get external IP",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c := ddns.New("y00ns.com")
		ip, err := c.GetExternalIP(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("external IP: %s", ip)
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
