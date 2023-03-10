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
	Short: "Use a configuration as current configuration",
	Long:  ``,
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
		err := kubeconfig.UseConfig(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("✅", color.GreenString(config))
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
