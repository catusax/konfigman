/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/catusax/konfigman/pkg/kubeconfig"
	"github.com/catusax/konfigman/pkg/tui"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a configuration",
	Long: `Use a configuration

soft link specific configuration to ~/.kube/config`,
	Run: func(cmd *cobra.Command, args []string) {
		var config string
		if len(args) == 0 {
			var err error
			config, err = tui.SelectConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else {
			config = args[0]
		}
		err := kubeconfig.UseConfig(config, useFlags)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("✅", color.GreenString(config))
	},
}

var useFlags kubeconfig.UseOptions

func loadUseFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&useFlags.Context, "context", "c", "", "context entry in kubeconfig")
	cmd.Flags().StringVarP(&useFlags.Namespace, "namespace", "n", "", "namespace for the context entry in kubeconfig")
	cmd.Flags().StringVar(&useFlags.Cluster, "cluster", "", "cluster for the context entry in kubeconfig")
	cmd.Flags().StringVarP(&useFlags.User, "user", "u", "", "user for the context entry in kubeconfig")
}

func init() {
	rootCmd.AddCommand(useCmd)
	loadUseFlags(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
