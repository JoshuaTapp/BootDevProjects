/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/JoshuaTapp/BootDevProjects/pokedexcli/internal/pokeAPI"
	"github.com/spf13/cobra"
)

var mapInit = false

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		mapNext()
	},
}

func init() {
	rootCmd.AddCommand(mapCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func mapNext() {
	locPtr := pokeAPI.GetLocations()

	if mapInit {
		err := locPtr.GetNext()

		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		mapInit = true
	}

	locPtr.PrintLocations()
}
