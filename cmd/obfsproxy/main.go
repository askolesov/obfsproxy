package main

import (
	"fmt"
	"os"

	"github.com/askolesov/obfsproxy/pkg"
	"github.com/spf13/cobra"
)

func main() {
	var listenAddr, targetAddr string

	rootCmd := &cobra.Command{
		Use:   "obfsproxy",
		Short: "A simple obfuscating proxy",
		Run: func(cmd *cobra.Command, args []string) {
			proxy := pkg.NewProxy(listenAddr, targetAddr)
			if err := proxy.Start(); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&listenAddr, "listen", "l", "localhost:8080", "Address to listen on")
	rootCmd.Flags().StringVarP(&targetAddr, "target", "t", "localhost:80", "Address to forward to")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
