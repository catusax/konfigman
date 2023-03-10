/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/catusax/konfigman/pkg/kubeconfig"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import file|ssh-url [name]",
	Short: "Import a kubeconfig file",
	Long: `You can import a kubeconfig from file or ssh

to import a kubeconfig from ssh, remote machine must have kubectl configured.
	`,
	Example: "kubectl kubeconfig import /path/to/kubeconfig",
	Args:    cobra.MatchAll(cobra.MaximumNArgs(1), cobra.MaximumNArgs(2)),
	Run: func(cmd *cobra.Command, args []string) {
		err := kubeconfig.InitRegistry()
		if err != nil {
			fmt.Println("failed to initialize:", err)
			os.Exit(1)
		}

		var name string
		var importoptions = make([]kubeconfig.ImportOption, 2)
		if len(args) >= 1 {
			url, err := url.Parse(args[0])
			if err != nil { //is filepath
				importoptions[0] = kubeconfig.ImportOptionWithPath(args[0])
				name = filepath.Base(args[0])
			} else {
				if url.Scheme == "ssh" {
					importoptions[0] = kubeconfig.ImportOptionWithSSH(args[0])
					url.Scheme = ""
					name = url.Host
				} else {
					fmt.Printf("bad path %s\n", args[0])
					os.Exit(1)
				}
			}
		} else {
			fmt.Printf("must provide an import source!\n")
			os.Exit(1)
		}

		if len(args) > 1 {
			name = args[1]
		}
		importoptions[1] = kubeconfig.ImportOptionWithName(name)

		err = kubeconfig.ImportConfig(importoptions...)
		if err != nil {
			if errors.Is(err, kubeconfig.ErrConfAlreadyExists) {
				fmt.Println("config already exists!")
			} else {
				fmt.Println("failed to import:", err)
				os.Exit(1)
			}

		} else {
			fmt.Println("imported", color.GreenString(name), "!")
		}

		if v, err := cmd.Flags().GetBool("use"); err == nil && v {
			err := kubeconfig.UseConfig(name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("✅", color.GreenString(name))
		}

	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().Bool("use", false, "use the imported configuration")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
