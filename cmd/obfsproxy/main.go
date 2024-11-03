package main

import (
	"fmt"
	"os"

	"github.com/askolesov/obfsproxy/pkg"
	"github.com/askolesov/obfsproxy/pkg/codec"
	"github.com/spf13/cobra"
)

func main() {
	var (
		listenAddr string
		targetAddr string
		isServer   bool
		isClient   bool
		key        string
		redundancy int
	)

	rootCmd := &cobra.Command{
		Use:   "obfsproxy",
		Short: "A simple obfuscating proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate flags
			if key == "" {
				return fmt.Errorf("key is required")
			}

			if isServer && isClient {
				return fmt.Errorf("cannot specify both server (-s) and client (-c) modes")
			}
			if !isServer && !isClient {
				isClient = true
			}

			if isClient {
				fmt.Println("Running in client mode")
			}
			if isServer {
				fmt.Println("Running in server mode")
			}

			// Calculate seed from key
			var seed uint64
			for _, ch := range key {
				seed += uint64(ch)
			}

			// Create codecs
			xorer, err := codec.NewXorer([]byte(key))
			if err != nil {
				return fmt.Errorf("failed to create xorer: %w", err)
			}

			injector, err := codec.NewInjector(seed, redundancy)
			if err != nil {
				return fmt.Errorf("failed to create injector: %w", err)
			}

			chain, err := codec.NewChain([]codec.Codec{
				codec.NewInverter(),
				xorer,
				injector,
			})
			if err != nil {
				return fmt.Errorf("failed to create codec chain: %w", err)
			}

			// Create proxy with codec
			proxy := pkg.NewProxy(listenAddr, targetAddr, isServer, chain)

			if err := proxy.Start(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	// Add flags
	rootCmd.Flags().StringVarP(&listenAddr, "listen", "l", "localhost:8080", "Address to listen on")
	rootCmd.Flags().StringVarP(&targetAddr, "target", "t", "localhost:80", "Address to forward to")
	rootCmd.Flags().BoolVarP(&isServer, "server", "s", false, "Run in server mode")
	rootCmd.Flags().BoolVarP(&isClient, "client", "c", false, "Run in client mode")
	rootCmd.Flags().StringVarP(&key, "key", "k", "", "Obfuscation key (required)")
	rootCmd.Flags().IntVarP(&redundancy, "redundancy", "r", 20, "Redundancy level (0-1000)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
