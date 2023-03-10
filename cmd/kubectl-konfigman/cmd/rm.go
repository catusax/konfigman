/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/catusax/konfigman/pkg/kubeconfig"
	"github.com/catusax/konfigman/pkg/tui"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		if len(args) > 0 {
			for _, arg := range args {
				err := kubeconfig.RemoveConfig(arg, force)
				if err != nil {
					if errors.Is(err, kubeconfig.ErrConfigInUse) {
						fmt.Println("config in use, use --force to force remove")
					} else {
						fmt.Println("remove config: ", err)
					}
				} else {
					fmt.Println(arg)
				}
			}
		} else {
			config, err := tui.SelectConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = kubeconfig.RemoveConfig(config, force)
			if err != nil {
				if errors.Is(err, kubeconfig.ErrConfigInUse) {
					fmt.Println("config in use, use --force to force remove")
				} else {
					fmt.Println("remove config: ", err)
				}
			} else {
				fmt.Println(config)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolP("force", "f", false, "force remove")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
