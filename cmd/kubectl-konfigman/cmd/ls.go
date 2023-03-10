/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/catusax/konfigman/pkg/kubeconfig"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all configurations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := kubeconfig.ListConfigs()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		current, _ := kubeconfig.GetCurrentConfig()

		for _, config := range configs {
			if current == config {
				fmt.Println("✅", color.GreenString(config))
			} else {

				fmt.Println("  ", config)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
